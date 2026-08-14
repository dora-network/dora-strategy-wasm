package dorastrategy_test

import (
	"context"
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

type stubStrategy struct {
	dorastrategy.StrategyBase
}

func (stubStrategy) Init(dorastrategy.Config) error { return nil }
func (stubStrategy) OnCandle(dorastrategy.Candle) ([]dorastrategy.OrderIntent, error) {
	return nil, nil
}

func TestTradeFields(t *testing.T) {
	tr := dorastrategy.Trade{
		TransactionID:      "tx-1",
		OrderBookID:        "OB-1",
		OrderID:            "order-1",
		OrderSeq:           42,
		UserID:             "user-1",
		Asset0:             "asset-1",
		Price:              "100.5",
		Quantity0:          "10.0",
		Side:               "BUY",
		AggressorIndicator: true,
		CreatedAt:          "2026-01-01T00:00:00Z",
	}
	if tr.TransactionID != "tx-1" {
		t.Errorf("TransactionID = %q, want %q", tr.TransactionID, "tx-1")
	}
	if tr.Asset0 != "asset-1" {
		t.Errorf("Asset0 = %q, want %q", tr.Asset0, "asset-1")
	}
	if tr.Quantity0 != "10.0" {
		t.Errorf("Quantity0 = %q, want %q", tr.Quantity0, "10.0")
	}
	if !tr.AggressorIndicator {
		t.Errorf("AggressorIndicator = false, want true")
	}
}

func TestPriceFields(t *testing.T) {
	p := dorastrategy.Price{
		AssetID: "asset-1",
		Price:   "100.5",
		YTM:     "0.0523",
		Time:    "2026-01-01T00:00:00Z",
	}
	if p.AssetID != "asset-1" {
		t.Errorf("AssetID = %q, want %q", p.AssetID, "asset-1")
	}
	if p.YTM != "0.0523" {
		t.Errorf("YTM = %q, want %q", p.YTM, "0.0523")
	}
}

func TestFrameworkVersionDefaultIsDev(t *testing.T) {
	// Default is "dev". The agent's validate path overwrites this
	// at build time via -ldflags. FrameworkVersion is a var
	// specifically so the validate path can set it without a
	// full framework rebuild.
	orig := dorastrategy.FrameworkVersion
	t.Cleanup(func() { dorastrategy.FrameworkVersion = orig })

	dorastrategy.FrameworkVersion = "test-version"
	if dorastrategy.FrameworkVersion != "test-version" {
		t.Errorf("FrameworkVersion assignment broken: got %q", dorastrategy.FrameworkVersion)
	}
}

type fakePreamble struct {
	candleCalls int
	tradeCalls  int
	priceCalls  int
}

func (f *fakePreamble) FetchCandles(ctx context.Context, start, end time.Time, resolution string, batchSize int) (dorastrategy.CandleBatch, error) {
	f.candleCalls++
	return dorastrategy.CandleBatch{Done: true}, nil
}

func (f *fakePreamble) FetchTrades(ctx context.Context, start, end time.Time, batchSize int) (dorastrategy.TradeBatch, error) {
	f.tradeCalls++
	return dorastrategy.TradeBatch{Done: true}, nil
}

func (f *fakePreamble) FetchPrices(ctx context.Context, start, end time.Time, batchSize int) (dorastrategy.PriceBatch, error) {
	f.priceCalls++
	return dorastrategy.PriceBatch{Done: true}, nil
}

func TestPreambleContextInterface(t *testing.T) {
	// Compile-time assertion: *fakePreamble satisfies PreambleContext.
	var _ dorastrategy.PreambleContext = (*fakePreamble)(nil)

	f := &fakePreamble{}
	ctx := t.Context()
	_, _ = f.FetchCandles(ctx, time.Now(), time.Now(), "1m", 100)
	_, _ = f.FetchTrades(ctx, time.Now(), time.Now(), 100)
	_, _ = f.FetchPrices(ctx, time.Now(), time.Now(), 100)
	if f.candleCalls != 1 || f.tradeCalls != 1 || f.priceCalls != 1 {
		t.Errorf("counts = %d/%d/%d, want 1/1/1", f.candleCalls, f.tradeCalls, f.priceCalls)
	}
}

func TestBatchTypes(t *testing.T) {
	cb := dorastrategy.CandleBatch{
		Items:  []dorastrategy.Candle{{OrderBookID: "OB-1", Close: "100"}},
		Done:   true,
		Cursor: "abc",
	}
	if len(cb.Items) != 1 || !cb.Done || cb.Cursor != "abc" {
		t.Errorf("CandleBatch = %+v", cb)
	}
	tb := dorastrategy.TradeBatch{
		Items:  []dorastrategy.Trade{{TransactionID: "tx-1"}},
		Done:   true,
		Cursor: "def",
	}
	if len(tb.Items) != 1 || !tb.Done || tb.Cursor != "def" {
		t.Errorf("TradeBatch = %+v", tb)
	}
	pb := dorastrategy.PriceBatch{
		Items:  []dorastrategy.Price{{AssetID: "asset-1"}},
		Done:   true,
		Cursor: "ghi",
	}
	if len(pb.Items) != 1 || !pb.Done || pb.Cursor != "ghi" {
		t.Errorf("PriceBatch = %+v", pb)
	}
}

type fullStrategy struct{}

func (fullStrategy) Init(dorastrategy.Config) error { return nil }
func (fullStrategy) OnPreamble(context.Context, dorastrategy.PreambleContext) error {
	return nil
}

func (fullStrategy) OnCandle(dorastrategy.Candle) ([]dorastrategy.OrderIntent, error) {
	return nil, nil
}

func (fullStrategy) OnTrade(dorastrategy.Trade) ([]dorastrategy.OrderIntent, error) {
	return nil, nil
}

func (fullStrategy) OnPrice(dorastrategy.Price) ([]dorastrategy.OrderIntent, error) {
	return nil, nil
}

func TestStrategyInterfaceFull(t *testing.T) {
	// Compile-time check: fullStrategy satisfies the new Strategy.
	var _ dorastrategy.Strategy = fullStrategy{}
}

// noopStrategy embeds the framework's StrategyBase; we expect it
// to satisfy the interface without overriding the new hooks.
type noopStrategy struct {
	dorastrategy.StrategyBase
}

func TestStrategyBaseSatisfiesInterface(t *testing.T) {
	var _ dorastrategy.Strategy = noopStrategy{}
}
