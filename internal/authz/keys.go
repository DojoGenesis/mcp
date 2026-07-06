// Package authz implements API-key authentication and dispatch authorization
// for the public HTTP mode. Keys are labeled so they are individually
// revocable and attributable in logs; key material itself is never logged.
package authz

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// minKeyLength is the minimum accepted key length. The operational standard
// is 32+ random bytes (hex/base64-encoded); this floor only guards against
// obvious misconfiguration like "label:test".
const minKeyLength = 16

var labelRe = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// KeySet holds the parsed label→key table.
type KeySet struct {
	keys map[string]string // label -> key
}

// ParseKeys parses DOJO_MCP_API_KEYS: comma-separated "label:key" pairs.
// It fails loudly on malformed entries so a misconfigured public endpoint
// refuses to start rather than serving with a half-loaded key table.
func ParseKeys(raw string) (*KeySet, error) {
	ks := &KeySet{keys: make(map[string]string)}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ks, nil
	}
	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		label, key, ok := strings.Cut(entry, ":")
		if !ok {
			return nil, fmt.Errorf("api key entry %q: want label:key", redactEntry(entry))
		}
		label = strings.TrimSpace(label)
		key = strings.TrimSpace(key)
		if !labelRe.MatchString(label) {
			return nil, fmt.Errorf("api key label %q: must match %s", label, labelRe)
		}
		if len(key) < minKeyLength {
			return nil, fmt.Errorf("api key %q: key shorter than %d chars", label, minKeyLength)
		}
		if strings.ContainsAny(key, " \t\r\n,") {
			return nil, fmt.Errorf("api key %q: key contains whitespace or comma", label)
		}
		if _, dup := ks.keys[label]; dup {
			return nil, fmt.Errorf("api key label %q: duplicate", label)
		}
		ks.keys[label] = key
	}
	return ks, nil
}

// Len returns the number of keys in the set.
func (ks *KeySet) Len() int { return len(ks.keys) }

// Labels returns the sorted key labels (safe to log).
func (ks *KeySet) Labels() []string {
	out := make([]string, 0, len(ks.keys))
	for l := range ks.keys {
		out = append(out, l)
	}
	sort.Strings(out)
	return out
}

// Authenticate returns the label owning the presented key. It compares
// against every key in the set without early exit, using SHA-256 digests and
// constant-time comparison so neither key length nor match position leaks.
func (ks *KeySet) Authenticate(presented string) (string, bool) {
	if presented == "" {
		return "", false
	}
	ph := sha256.Sum256([]byte(presented))
	var matched string
	found := 0
	for label, key := range ks.keys {
		kh := sha256.Sum256([]byte(key))
		if subtle.ConstantTimeCompare(ph[:], kh[:]) == 1 {
			matched = label
			found = 1
		}
	}
	return matched, found == 1
}

// redactEntry keeps only the label half of a malformed entry for error text.
func redactEntry(entry string) string {
	if label, _, ok := strings.Cut(entry, ":"); ok {
		return label + ":***"
	}
	if len(entry) > 8 {
		return entry[:8] + "***"
	}
	return entry
}
