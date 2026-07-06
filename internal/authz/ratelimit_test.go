package authz

import (
	"context"
	"testing"
	"time"
)

func TestLimiter_BurstThenBlock(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	l := NewLimiter(3)
	l.now = func() time.Time { return now }

	for i := 0; i < 3; i++ {
		if !l.Allow("win") {
			t.Fatalf("call %d: want allowed (burst)", i+1)
		}
	}
	if l.Allow("win") {
		t.Fatal("4th immediate call: want blocked")
	}
}

func TestLimiter_RefillsOverTime(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	l := NewLimiter(6) // one token per 10s
	l.now = func() time.Time { return now }

	for i := 0; i < 6; i++ {
		l.Allow("win")
	}
	if l.Allow("win") {
		t.Fatal("bucket should be empty")
	}

	now = now.Add(10 * time.Second)
	if !l.Allow("win") {
		t.Fatal("one token should have refilled after 10s at 6/min")
	}
	if l.Allow("win") {
		t.Fatal("only one token should have refilled")
	}
}

func TestLimiter_LabelsIndependent(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	l := NewLimiter(1)
	l.now = func() time.Time { return now }

	if !l.Allow("win") {
		t.Fatal("win first call blocked")
	}
	if l.Allow("win") {
		t.Fatal("win second call allowed")
	}
	if !l.Allow("mac") {
		t.Fatal("mac should have its own bucket")
	}
}

func TestLimiter_NilAllowsEverything(t *testing.T) {
	var l *Limiter
	if !l.Allow("anyone") {
		t.Fatal("nil limiter must allow")
	}
}

func TestDispatchAllowed_DefaultsTrueWithoutIdentity(t *testing.T) {
	ctx := context.Background()
	if !DispatchAllowed(ctx) {
		t.Fatal("no identity (stdio mode) must allow dispatch")
	}
	if _, ok := Label(ctx); ok {
		t.Fatal("no identity should have no label")
	}

	ctx = WithIdentity(ctx, "web", false)
	if DispatchAllowed(ctx) {
		t.Fatal("identity with dispatchAllowed=false must deny")
	}
	if label, ok := Label(ctx); !ok || label != "web" {
		t.Fatalf("Label = %q,%v; want web,true", label, ok)
	}

	ctx = WithIdentity(ctx, "ops", true)
	if !DispatchAllowed(ctx) {
		t.Fatal("identity with dispatchAllowed=true must allow")
	}
}
