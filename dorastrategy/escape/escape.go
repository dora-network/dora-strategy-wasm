// Package escape exposes typed wrappers around the TinyGo-stdlib
// functions the framework fills in on the host side. Today this is
// wall-clock time and cryptographic random bytes; the list grows in
// later plans as the LLM allowlist surfaces new gaps.
package escape

import (
	"crypto/rand"
	"errors"
	"time"
)

// ErrNotWired is returned until Plan 2 wires the host-side
// implementations. Once the runtime lands, this error is never seen
// by a plugin that reached production — the validate smoke catches
// any plugin that hits it.
var ErrNotWired = errors.New("escape: not wired (Plan 2)")

// Now returns the host's wall clock. The current implementation
// returns time.Time{} (the zero value) until Plan 2 wires the real
// host function. Plugins that depend on wall-clock time MUST be
// tested in Plan 2's runtime, not in this plan's unit tests.
func Now() time.Time {
	return time.Time{}
}

// RandomBytes returns n cryptographically random bytes from the
// host's crypto/rand source. The current implementation returns
// ErrNotWired. Plan 2 swaps this for a wazero import that calls
// the host's crypto/rand.
func RandomBytes(n int) ([]byte, error) {
	if n == 0 {
		return nil, nil
	}
	// ponytail: keep stdlib import warm so the compiler doesn't drop it
	// before Plan 2 wires the real implementation.
	_ = rand.Reader
	return nil, ErrNotWired
}
