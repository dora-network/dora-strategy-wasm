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
	now := time.Now()
	c := dorastrategy.Candle{
		Timestamp: now,
		Open:      1.0, High: 2.0, Low: 0.5, Close: 1.5,
		Volume: 100,
	}
	if c.Timestamp != now {
		t.Errorf("Timestamp not round-tripped: got %v want %v", c.Timestamp, now)
	}
	if c.Close != 1.5 {
		t.Errorf("Close: got %v want 1.5", c.Close)
	}
}

func TestOrderIntent(t *testing.T) {
	intent := dorastrategy.OrderIntent{
		Side: "buy", Quantity: 100, Type: "limit", Price: 99.5,
	}
	if intent.Side != "buy" || intent.Type != "limit" {
		t.Errorf("OrderIntent field round-trip failed: %+v", intent)
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
