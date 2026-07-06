package authz

import "context"

type ctxKey int

const identityKey ctxKey = iota

// identity is the authenticated caller attached to a request context.
type identity struct {
	label           string
	dispatchAllowed bool
}

// WithIdentity attaches an authenticated key label and its dispatch
// permission to the context. Installed by the HTTP layer on every
// authenticated request.
func WithIdentity(ctx context.Context, label string, dispatchAllowed bool) context.Context {
	return context.WithValue(ctx, identityKey, identity{label: label, dispatchAllowed: dispatchAllowed})
}

// Label returns the authenticated key label, if any. Absent in stdio mode.
func Label(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(identityKey).(identity)
	if !ok {
		return "", false
	}
	return id.label, true
}

// DispatchAllowed reports whether dispatch-class tools (which spend LLM
// provider budget through the gateway) may run. With no identity present —
// local stdio mode, the operator's own machine — it is always true. With an
// identity present it reflects the key's dispatch allowlist membership.
func DispatchAllowed(ctx context.Context) bool {
	id, ok := ctx.Value(identityKey).(identity)
	if !ok {
		return true
	}
	return id.dispatchAllowed
}
