package escape_test

import (
	"errors"
	"testing"
	"time"

	"github.com/dora-network/dora-strategy-wasm/dorastrategy/escape"
)

func TestNow_NotWired_ReturnsZero(t *testing.T) {
	got := escape.Now()
	if !got.IsZero() {
		t.Errorf("Now returned %v, want zero time", got)
	}
}

func TestNow_UnwiredSentinel(t *testing.T) {
	// Plan 2 swaps this for the real host call. Until then, the value
	// is the zero time and the plugin's only contract is that the
	// returned time is not "random".
	got := escape.Now()
	if got.Location() != time.UTC && got.Location().String() != "" {
		// zero time has Location() == time.UTC, which is fine.
		t.Errorf("Now returned location %v, want zero", got.Location())
	}
}

func TestRandomBytes_NotWired(t *testing.T) {
	_, err := escape.RandomBytes(16)
	if !errors.Is(err, escape.ErrNotWired) {
		t.Fatalf("RandomBytes returned err=%v, want ErrNotWired", err)
	}
}

func TestRandomBytes_ZeroLength(t *testing.T) {
	// Zero-length calls return an empty slice and ErrNotWired (the
	// host wiring in Plan 2 short-circuits on zero-length). Either
	// behavior is acceptable; assert the function does not panic.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("RandomBytes(0) panicked: %v", r)
		}
	}()
	_, _ = escape.RandomBytes(0)
}
