package dorastrategy_test

import (
	"testing"
	"time"

	"github.com/dora-network/dora-strategy-wasm/dorastrategy"
)

func TestModeConstants(t *testing.T) {
	if dorastrategy.ModeValidate == "" {
		t.Fatal("ModeValidate must be non-empty")
	}
	if dorastrategy.ModeBacktest == "" {
		t.Fatal("ModeBacktest must be non-empty")
	}
	if dorastrategy.ModeLive == "" {
		t.Fatal("ModeLive must be non-empty")
	}
}

func TestCandleFields(t *testing.T) {
	c := dorastrategy.Candle{
		OrderBookID:    "OB-1234",
		StartTimestamp: "2026-01-01T00:00:00Z",
		Open:           "100.5",
		High:           "101.2",
		Low:            "99.8",
		Close:          "100.9",
		OpenYtm:        "0.0523",
		CloseYtm:       "0.0519",
		HighYtm:        "0.0525",
		LowYtm:         "0.0517",
		Volume:         "12500",
	}
	if c.OrderBookID != "OB-1234" {
		t.Errorf("OrderBookID: got %q", c.OrderBookID)
	}
	if c.Open != "100.5" {
		t.Errorf("Open: got %q want decimal string", c.Open)
	}
	if c.OpenYtm != "0.0523" {
		t.Errorf("OpenYtm: got %q", c.OpenYtm)
	}
}

func TestOrderIntent(t *testing.T) {
	intent := dorastrategy.OrderIntent{
		Side:               "buy",
		Quantity:           "100",
		Type:               "limit",
		Price:              "99.5",
		InverseLeverage:    "2",
		FromGlobalPosition: false,
	}
	if intent.Side != "buy" || intent.Type != "limit" {
		t.Errorf("OrderIntent field round-trip failed: %+v", intent)
	}
	if intent.Quantity != "100" || intent.Price != "99.5" {
		t.Errorf("decimal-string fields round-trip: %+v", intent)
	}
	if intent.InverseLeverage != "2" {
		t.Errorf("InverseLeverage: got %q want \"2\"", intent.InverseLeverage)
	}
	if intent.FromGlobalPosition {
		t.Errorf("FromGlobalPosition: got true want false (default isolated margin)")
	}
}

// TestOrderIntent_Defaults verifies the zero value semantics: a fresh
// OrderIntent with all fields unset represents a buy-market order
// with no leverage on the isolated margin account. The framework
// applies these defaults in the host; the plugin's code may rely on
// them by leaving the fields blank.
func TestOrderIntent_Defaults(t *testing.T) {
	var intent dorastrategy.OrderIntent
	if intent.Side != "" || intent.Quantity != "" || intent.Type != "" {
		t.Errorf("zero value should be all-empty: %+v", intent)
	}
	if intent.FromGlobalPosition {
		t.Errorf("zero-value FromGlobalPosition must be false (isolated margin)")
	}
	if intent.InverseLeverage != "" {
		t.Errorf("zero-value InverseLeverage must be empty (host default \"1\")")
	}
}

func TestConfigRoundTrip(t *testing.T) {
	cfg := dorastrategy.Config{
		Mode:        dorastrategy.ModeBacktest,
		OrderBookID: "OB-1234",
		Start:       time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		End:         time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		Resolution:  "5m",
		DoraBaseURL: "https://staging.dora.co",
		DoraAPIKey:  "key",
		Params:      map[string]string{"stop_loss_bps": "50"},
	}
	if cfg.Mode != dorastrategy.ModeBacktest {
		t.Errorf("Mode: got %v want %v", cfg.Mode, dorastrategy.ModeBacktest)
	}
	if cfg.Params["stop_loss_bps"] != "50" {
		t.Errorf("Params: got %v want stop_loss_bps=50", cfg.Params)
	}
}

func TestStrategyInterface(t *testing.T) {
	// Compile-time assertion that a type implementing the two methods
	// satisfies dorastrategy.Strategy.
	var _ dorastrategy.Strategy = stubStrategy{}
}

type stubStrategy struct{}

func (stubStrategy) Init(dorastrategy.Config) error { return nil }
func (stubStrategy) OnCandle(dorastrategy.Candle) ([]dorastrategy.OrderIntent, error) {
	return nil, nil
}
