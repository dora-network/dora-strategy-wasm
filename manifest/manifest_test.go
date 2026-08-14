package manifest_test

import (
	"strings"
	"testing"

	"github.com/dora-network/dora-strategy-wasm/dorastrategy"
	"github.com/dora-network/dora-strategy-wasm/manifest"
)

func validManifest() manifest.Manifest {
	return manifest.Manifest{
		SchemaVersion:    1,
		ModuleName:       "vwap-momentum",
		FrameworkVersion: dorastrategy.FrameworkVersion,
		Language:         "go",
		Capabilities: manifest.Capabilities{
			OrderBooks:    []string{"OB-1234"},
			Resolutions:   []string{"5m"},
			Channels:      []string{"candle", "trade"},
			HostFunctions: []string{"host_log", "host_submit_order"},
		},
		ParamsSchema: map[string]string{"stop_loss_bps": "int"},
		Preamble:     manifest.PreambleCapabilities{WarmupCandles: 0},
	}
}

func TestRoundTrip(t *testing.T) {
	m := validManifest()
	b, err := m.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	var got manifest.Manifest
	if err := got.UnmarshalJSON(b); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}
	if got.ModuleName != m.ModuleName {
		t.Errorf("ModuleName: got %q want %q", got.ModuleName, m.ModuleName)
	}
	if len(got.Capabilities.OrderBooks) != 1 || got.Capabilities.OrderBooks[0] != "OB-1234" {
		t.Errorf("OrderBooks round-trip: %v", got.Capabilities.OrderBooks)
	}
}

func TestValidate_Happy(t *testing.T) {
	if err := validManifest().Validate(); err != nil {
		t.Errorf("valid manifest rejected: %v", err)
	}
}

func TestValidate_BadSchemaVersion(t *testing.T) {
	m := validManifest()
	m.SchemaVersion = 99
	if err := m.Validate(); err == nil || !strings.Contains(err.Error(), "schema_version") {
		t.Fatalf("expected schema_version error, got %v", err)
	}
}

func TestValidate_EmptyOrderBooks(t *testing.T) {
	m := validManifest()
	m.Capabilities.OrderBooks = nil
	if err := m.Validate(); err == nil || !strings.Contains(err.Error(), "order_books") {
		t.Fatalf("expected order_books error, got %v", err)
	}
}

func TestValidate_EmptyHostFunctions(t *testing.T) {
	m := validManifest()
	m.Capabilities.HostFunctions = nil
	if err := m.Validate(); err == nil || !strings.Contains(err.Error(), "host_functions") {
		t.Fatalf("expected host_functions error, got %v", err)
	}
}

func TestValidate_UnknownParamType(t *testing.T) {
	m := validManifest()
	m.ParamsSchema["weird"] = "datetime"
	if err := m.Validate(); err == nil || !strings.Contains(err.Error(), "params_schema") {
		t.Fatalf("expected params_schema error, got %v", err)
	}
}

func TestManifestRequiresFrameworkVersion(t *testing.T) {
	m := manifest.Manifest{
		SchemaVersion: 1,
		ModuleName:    "m",
		Capabilities:  manifest.Capabilities{OrderBooks: []string{"OB-1"}, Resolutions: []string{"1m"}, Channels: []string{"candle"}, HostFunctions: []string{"host_log"}},
		Preamble:      manifest.PreambleCapabilities{},
	}
	if err := m.Validate(); err == nil {
		t.Fatal("expected error for missing framework_version")
	}
}

func TestManifestFrameworkVersionMismatch(t *testing.T) {
	m := manifest.Manifest{
		SchemaVersion:    1,
		ModuleName:       "m",
		FrameworkVersion: "v0",
		Capabilities:     manifest.Capabilities{OrderBooks: []string{"OB-1"}, Resolutions: []string{"1m"}, Channels: []string{"candle"}, HostFunctions: []string{"host_log"}},
		Preamble:         manifest.PreambleCapabilities{},
	}
	if err := m.Validate(); err == nil {
		t.Fatal("expected error for framework_version mismatch")
	}
}

func TestManifestPriceChannelRequiresAssetID(t *testing.T) {
	m := manifest.Manifest{
		SchemaVersion:    1,
		ModuleName:       "m",
		FrameworkVersion: dorastrategy.FrameworkVersion,
		Capabilities:     manifest.Capabilities{OrderBooks: []string{"OB-1"}, Resolutions: []string{"1m"}, Channels: []string{"candle", "price"}, HostFunctions: []string{"host_log"}},
		Preamble:         manifest.PreambleCapabilities{},
	}
	if err := m.Validate(); err == nil {
		t.Fatal("expected error for channels=[price] without asset_ids")
	}
}

func TestManifestPriceChannelRequiresExactlyOneAssetID(t *testing.T) {
	m := manifest.Manifest{
		SchemaVersion:    1,
		ModuleName:       "m",
		FrameworkVersion: dorastrategy.FrameworkVersion,
		Capabilities:     manifest.Capabilities{OrderBooks: []string{"OB-1"}, Resolutions: []string{"1m"}, Channels: []string{"candle", "price"}, AssetIDs: []string{"a1", "a2"}, HostFunctions: []string{"host_log"}},
		Preamble:         manifest.PreambleCapabilities{},
	}
	if err := m.Validate(); err == nil {
		t.Fatal("expected error for channels=[price] with 2 asset_ids")
	}
}
