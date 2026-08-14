// Package manifest defines the capability declaration a generated
// WASM plugin emits alongside its .wasm blob. The host validates the
// manifest at load time: a plugin that calls a host function it did
// not declare is rejected at instantiate time.
package manifest

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dora-network/dora-strategy-wasm/dorastrategy"
)

// SupportedSchemaVersion is the only schema_version the host
// understands. A future bump adds a new constant and a migration.
const SupportedSchemaVersion = 1

// AllowedParamTypes is the closed set of param types a manifest may
// declare. Anything outside this set fails Validate().
//
//nolint:gochecknoglobals // ponytail: constant set, not state.
var AllowedParamTypes = map[string]struct{}{
	"string": {},
	"int":    {},
	"float":  {},
	"bool":   {},
}

// Manifest is the JSON object a plugin emits next to its .wasm.
type Manifest struct {
	SchemaVersion    int    `json:"schema_version"`
	ModuleName       string `json:"module_name"`
	FrameworkVersion string `json:"framework_version"`
	Language         string `json:"language"`
	//nolint:tagliatelle // standard TinyGo project name; rename would diverge from spec.
	TinyGoVersion string               `json:"tinygo_version,omitempty"`
	GoVersion     string               `json:"go_version,omitempty"`
	Capabilities  Capabilities         `json:"capabilities"`
	ParamsSchema  map[string]string    `json:"params_schema"`
	Preamble      PreambleCapabilities `json:"preamble"`
}

// Capabilities declares what the plugin is allowed to do.
type Capabilities struct {
	OrderBooks    []string `json:"order_books"`
	Resolutions   []string `json:"resolutions,omitempty"`
	Channels      []string `json:"channels,omitempty"`  // "candle" | "trade" | "price"
	AssetIDs      []string `json:"asset_ids,omitempty"` //nolint:tagliatelle // asyncapi name; required when "price" in Channels
	HostFunctions []string `json:"host_functions"`
}

// PreambleCapabilities declares the constraints on the strategy's
// OnPreamble phase. The preamble is read-only historical analysis;
// the host gates order submission.
type PreambleCapabilities struct {
	WarmupCandles int `json:"warmup_candles"` // resolution-spaced candles
}

// Validate enforces the spec's mandatory fields. Returns nil on
// success, a typed error describing the first failure otherwise.
func (m Manifest) Validate() error {
	if m.FrameworkVersion == "" {
		return errors.New("manifest: framework_version is required")
	}
	if m.FrameworkVersion != dorastrategy.FrameworkVersion {
		return fmt.Errorf("manifest: framework_version %q does not match expected %q (rebuild the strategy against the current framework)",
			m.FrameworkVersion, dorastrategy.FrameworkVersion)
	}
	if m.Preamble.WarmupCandles < 0 {
		return errors.New("manifest: preamble.warmup_candles must be >= 0")
	}
	hasPrice := false
	for _, ch := range m.Capabilities.Channels {
		if ch == "price" {
			hasPrice = true
			break
		}
	}
	if hasPrice && len(m.Capabilities.AssetIDs) != 1 {
		return fmt.Errorf("manifest: channels includes 'price' but asset_ids has %d entries "+
			"(must be exactly 1 for the POC)", len(m.Capabilities.AssetIDs))
	}
	if m.SchemaVersion != SupportedSchemaVersion {
		return fmt.Errorf("manifest: schema_version %d not supported (want %d)",
			m.SchemaVersion, SupportedSchemaVersion)
	}
	if m.ModuleName == "" {
		return errors.New("manifest: module_name is required")
	}
	if len(m.Capabilities.OrderBooks) == 0 {
		return errors.New("manifest: capabilities.order_books must be non-empty")
	}
	if len(m.Capabilities.HostFunctions) == 0 {
		return errors.New("manifest: capabilities.host_functions must be non-empty")
	}
	for k, v := range m.ParamsSchema {
		if _, ok := AllowedParamTypes[v]; !ok {
			return fmt.Errorf("manifest: params_schema.%s has type %q (not in %v)",
				k, v, keys(AllowedParamTypes))
		}
	}
	return nil
}

// MarshalJSON encodes the manifest as JSON.
func (m Manifest) MarshalJSON() ([]byte, error) {
	type alias Manifest
	return json.Marshal(alias(m))
}

// UnmarshalJSON decodes JSON into the manifest.
func (m *Manifest) UnmarshalJSON(b []byte) error {
	type alias Manifest
	return json.Unmarshal(b, (*alias)(m))
}

func keys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
