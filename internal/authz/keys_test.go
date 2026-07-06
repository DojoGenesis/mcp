package authz

import (
	"reflect"
	"strings"
	"testing"
)

const (
	testKeyA = "0123456789abcdef0123456789abcdef"
	testKeyB = "fedcba9876543210fedcba9876543210"
)

func TestParseKeys_Valid(t *testing.T) {
	ks, err := ParseKeys("win:" + testKeyA + " , mac:" + testKeyB)
	if err != nil {
		t.Fatalf("ParseKeys: %v", err)
	}
	if ks.Len() != 2 {
		t.Fatalf("Len = %d, want 2", ks.Len())
	}
	if got := ks.Labels(); !reflect.DeepEqual(got, []string{"mac", "win"}) {
		t.Fatalf("Labels = %v, want [mac win]", got)
	}
}

func TestParseKeys_Empty(t *testing.T) {
	ks, err := ParseKeys("  ")
	if err != nil {
		t.Fatalf("ParseKeys empty: %v", err)
	}
	if ks.Len() != 0 {
		t.Fatalf("Len = %d, want 0", ks.Len())
	}
	if _, ok := ks.Authenticate(testKeyA); ok {
		t.Fatal("empty KeySet authenticated a key")
	}
}

func TestParseKeys_Malformed(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string // substring of the error
	}{
		{"no colon", "justakeywithoutlabel0000000000", "want label:key"},
		{"short key", "win:short", "shorter than"},
		{"bad label", "win win:" + testKeyA, "must match"},
		{"duplicate label", "win:" + testKeyA + ",win:" + testKeyB, "duplicate"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseKeys(tc.raw)
			if err == nil {
				t.Fatalf("ParseKeys(%q): want error, got nil", tc.raw)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error %q does not contain %q", err, tc.want)
			}
		})
	}
}

func TestParseKeys_ErrorNeverContainsKey(t *testing.T) {
	_, err := ParseKeys("bad label:" + testKeyA)
	if err == nil {
		t.Fatal("want error")
	}
	if strings.Contains(err.Error(), testKeyA) {
		t.Fatalf("error leaks key material: %v", err)
	}
}

func TestAuthenticate(t *testing.T) {
	ks, err := ParseKeys("win:" + testKeyA + ",mac:" + testKeyB)
	if err != nil {
		t.Fatalf("ParseKeys: %v", err)
	}

	if label, ok := ks.Authenticate(testKeyA); !ok || label != "win" {
		t.Fatalf("Authenticate(keyA) = %q,%v; want win,true", label, ok)
	}
	if label, ok := ks.Authenticate(testKeyB); !ok || label != "mac" {
		t.Fatalf("Authenticate(keyB) = %q,%v; want mac,true", label, ok)
	}
	if _, ok := ks.Authenticate("wrong-key-material-000000000000"); ok {
		t.Fatal("wrong key authenticated")
	}
	if _, ok := ks.Authenticate(""); ok {
		t.Fatal("empty key authenticated")
	}
	// Prefix of a real key must not match.
	if _, ok := ks.Authenticate(testKeyA[:20]); ok {
		t.Fatal("key prefix authenticated")
	}
}
