package allowlist_test

import (
	"testing"

	"github.com/dora-network/dora-strategy-wasm/allowlist"
)

func TestIsAllowed_FrameworkImports(t *testing.T) {
	allowed := []string{
		"github.com/dora-network/dora-strategy-wasm/dorastrategy",
		"github.com/dora-network/dora-strategy-wasm/dorastrategy/host",
		"github.com/dora-network/dora-strategy-wasm/dorastrategy/escape",
		"github.com/dora-network/dora-strategy-wasm/manifest",
	}
	for _, p := range allowed {
		if !allowlist.IsAllowed(p) {
			t.Errorf("IsAllowed(%q) = false, want true", p)
		}
	}
}

func TestIsAllowed_StdlibSafeSubset(t *testing.T) {
	safe := []string{
		"fmt", "errors", "strings", "strconv", "time", "math", "sort",
		"context", "sync", "encoding/json", "io",
	}
	for _, p := range safe {
		if !allowlist.IsAllowed(p) {
			t.Errorf("IsAllowed(%q) = false, want true (stdlib safe subset)", p)
		}
	}
}

func TestIsAllowed_StdlibDenied(t *testing.T) {
	denied := []string{
		"net/http", "net", "os", "os/exec", "database/sql",
		"plugin", "reflect", "unsafe", "syscall", "crypto/tls",
		"runtime/debug",
	}
	for _, p := range denied {
		if allowlist.IsAllowed(p) {
			t.Errorf("IsAllowed(%q) = true, want false", p)
		}
	}
}

func TestIsAllowed_DoraSDKDenied(t *testing.T) {
	denied := []string{
		"github.com/dora-network/dora-client-go",
		"github.com/dora-network/dora-client-go/doraclient",
		"github.com/coder/websocket",
	}
	for _, p := range denied {
		if allowlist.IsAllowed(p) {
			t.Errorf("IsAllowed(%q) = true, want false (Dora SDK / wsplex)", p)
		}
	}
}

func TestIsAllowed_UnknownImport(t *testing.T) {
	if allowlist.IsAllowed("github.com/some/random/thirdparty") {
		t.Error("IsAllowed(third-party) = true, want false (default-deny)")
	}
}
