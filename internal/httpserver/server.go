// Package httpserver wires the MCP streamable-HTTP transport behind
// Bearer-key authentication for the public endpoint (Lane B). Ingress is
// expected to be a Cloudflare tunnel; this server still authenticates every
// /mcp request itself and never logs key material.
package httpserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/DojoGenesis/mcp/internal/authz"
)

// Options configures the public HTTP handler.
type Options struct {
	Version        string
	Keys           *authz.KeySet
	DispatchLabels map[string]bool // key labels allowed to run dispatch-class tools
	MCP            http.Handler    // the mcp-go streamable HTTP server
}

// NewHandler builds the route table:
//
//	GET  /health   — unauthenticated liveness (status+version only)
//	*    /mcp      — Bearer-key auth (Authorization header)
//	*    /mcp/k/…  — key-in-path fallback for clients that cannot send
//	                 custom headers; same key table, path is redacted in logs
func NewHandler(opts Options) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "dojo-mcp-server",
			"version": opts.Version,
		})
	})

	mux.Handle("/mcp", opts.bearerAuth(opts.MCP))
	mux.Handle("/mcp/k/", opts.pathKeyAuth(opts.MCP))

	return requestLog(mux)
}

// IdentityContextFunc is the mcp-go HTTPContextFunc bridge: it copies the
// authenticated identity from the HTTP request context (set by the auth
// middleware) into the context handed to tool handlers.
func IdentityContextFunc(ctx context.Context, r *http.Request) context.Context {
	if label, ok := authz.Label(r.Context()); ok {
		return authz.WithIdentity(ctx, label, authz.DispatchAllowed(r.Context()))
	}
	return ctx
}

// Run serves h on addr with timeouts safe for SSE streaming (no write
// timeout — streams are long-lived) and shuts down gracefully on
// SIGINT/SIGTERM.
func Run(addr string, h http.Handler) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case s := <-sig:
		log.Printf("httpserver: received %v, shutting down", s)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	}
}

// ─── auth middlewares ─────────────────────────────────────────────────────────

func (o Options) bearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r.Header.Get("Authorization"))
		label, ok := o.Keys.Authenticate(token)
		if !ok {
			unauthorized(w, r)
			return
		}
		recordLabel(r, label)
		ctx := authz.WithIdentity(r.Context(), label, o.DispatchLabels[label])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// pathKeyAuth serves /mcp/k/<key>[/...] for clients that cannot set an
// Authorization header. The key segment is authenticated against the same
// table, then the path is rewritten to /mcp before reaching the transport.
func (o Options) pathKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rest := strings.TrimPrefix(r.URL.Path, "/mcp/k/")
		key := rest
		var tail string
		if i := strings.IndexByte(rest, '/'); i >= 0 {
			key, tail = rest[:i], rest[i:]
		}
		label, ok := o.Keys.Authenticate(key)
		if !ok {
			unauthorized(w, r)
			return
		}
		recordLabel(r, label)
		r2 := r.Clone(authz.WithIdentity(r.Context(), label, o.DispatchLabels[label]))
		r2.URL.Path = "/mcp" + tail
		r2.URL.RawPath = ""
		next.ServeHTTP(w, r2)
	})
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return strings.TrimSpace(header[len(prefix):])
	}
	return ""
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	log.Printf("http auth_fail remote=%s path=%s", r.RemoteAddr, sanitizePath(r.URL.Path))
	w.Header().Set("WWW-Authenticate", `Bearer realm="dojo-mcp"`)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
}

// ─── request logging ──────────────────────────────────────────────────────────

// labelHolder lets the auth middleware report the authenticated label back
// up to the request logger (context values only flow downward).
type labelHolder struct{ label string }

type holderCtxKey struct{}

func recordLabel(r *http.Request, label string) {
	if h, ok := r.Context().Value(holderCtxKey{}).(*labelHolder); ok {
		h.label = label
	}
}

// statusRecorder captures the response status while passing SSE flushes
// through to the underlying writer.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Flush() {
	if f, ok := sr.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		holder := &labelHolder{}
		r = r.WithContext(context.WithValue(r.Context(), holderCtxKey{}, holder))
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		label := holder.label
		if label == "" {
			label = "-"
		}
		log.Printf("http %s %s status=%d key=%s dur_ms=%d",
			r.Method, sanitizePath(r.URL.Path), rec.status, label, time.Since(start).Milliseconds())
	})
}

// sanitizePath redacts key material embedded in /mcp/k/ paths.
func sanitizePath(p string) string {
	if strings.HasPrefix(p, "/mcp/k/") {
		return "/mcp/k/***"
	}
	return p
}
