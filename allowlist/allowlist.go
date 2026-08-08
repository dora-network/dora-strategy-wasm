// Package allowlist enumerates the imports a generated WASM plugin
// may use. The validator calls IsAllowed for every import in the
// LLM-authored go.mod; an import outside the allow set fails the
// build and the diagnostic is fed back to the model for repair.
package allowlist

import "strings"

// allowedStdlib is the TinyGo-compatible subset of the standard
// library. Anything not in this set is denied.
//
//nolint:gochecknoglobals // ponytail: constant data, not state.
var allowedStdlib = map[string]struct{}{
	"fmt": {}, "errors": {}, "strings": {}, "strconv": {},
	"time": {}, "math": {}, "sort": {}, "context": {}, "sync": {},
	"encoding/json": {}, "io": {}, "bytes": {}, "bufio": {},
	"unicode": {}, "unicode/utf8": {}, "unicode/utf16": {},
}

// denied is checked first. Anything here is rejected regardless of
// any allow. The deny list is intentionally explicit so a future
// audit can see exactly which imports the policy forbids.
//
//nolint:gochecknoglobals // ponytail: constant data, not state.
var denied = map[string]struct{}{
	"net/http":                               {},
	"net":                                    {},
	"os":                                     {},
	"os/exec":                                {},
	"os/signal":                              {},
	"database/sql":                           {},
	"plugin":                                 {},
	"reflect":                                {},
	"unsafe":                                 {},
	"syscall":                                {},
	"crypto/tls":                             {},
	"runtime/debug":                          {},
	"runtime/pprof":                          {},
	"runtime/trace":                          {},
	"github.com/dora-network/dora-client-go": {},
	"github.com/dora-network/dora-client-go/doraclient": {},
	"github.com/coder/websocket":                        {},
}

// allowedFramework lists every package in this sub-module the
// strategy may import. The keys are import path prefixes; a plugin
// that imports a sub-package (e.g. .../dorastrategy/host) matches
// the parent prefix.
//
//nolint:gochecknoglobals // ponytail: constant prefix list, not state.
var allowedFramework = []string{
	"github.com/dora-network/dora-strategy-wasm/dorastrategy",
	"github.com/dora-network/dora-strategy-wasm/manifest",
	"github.com/dora-network/dora-strategy-wasm/allowlist",
}

// IsAllowed returns true iff the import path is in the allow set and
// not in the deny set. Default-deny: an unknown import is rejected.
//
// The check order is deny-first, then framework-prefix, then
// stdlib-exact, then default-deny. A future plan may add a curated
// set of third-party imports (e.g. a stats library); the check
// order makes that addition a one-line change.
func IsAllowed(importPath string) bool {
	if _, bad := denied[importPath]; bad {
		return false
	}
	for _, prefix := range allowedFramework {
		if strings.HasPrefix(importPath, prefix) {
			return true
		}
	}
	if _, ok := allowedStdlib[importPath]; ok {
		return true
	}
	return false
}
