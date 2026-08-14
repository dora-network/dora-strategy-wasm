package host_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/dora-network/dora-strategy-wasm/dorastrategy/host"
)

func TestLevelConstants(t *testing.T) {
	if host.LevelDebug == "" || host.LevelInfo == "" ||
		host.LevelWarn == "" || host.LevelError == "" {
		t.Fatal("all four Level constants must be non-empty")
	}
}

func TestCandleMirrorTags(t *testing.T) {
	c := host.Candle{
		OrderBookID:    "OB-1",
		StartTimestamp: "2026-01-01T00:00:00Z",
		Open:           "1", High: "2", Low: "3", Close: "4",
		OpenYtm: "0.1", CloseYtm: "0.2", HighYtm: "0.3", LowYtm: "0.4",
		Volume: "100",
	}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	want := []string{
		"order_book_id", "start_timestamp", "open", "high", "low", "close",
		"open_ytm", "close_ytm", "high_ytm", "low_ytm", "volume",
	}
	for _, k := range want {
		if _, ok := m[k]; !ok {
			t.Errorf("missing JSON key %q in %s", k, string(b))
		}
	}
}

func TestFillMirrorTags(t *testing.T) {
	f := host.Fill{OrderID: "o1", Price: "1.5", Quantity: "2", Simulated: true}
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, k := range []string{"order_id", "price", "quantity", "simulated"} {
		if _, ok := m[k]; !ok {
			t.Errorf("missing JSON key %q in %s", k, string(b))
		}
	}
}

func TestConfigMirrorTags(t *testing.T) {
	cfg := host.Config{
		Mode:        "backtest",
		OrderBookID: "OB-1",
		Start:       time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		End:         time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		Resolution:  "5m",
		DoraBaseURL: "https://dora.co",
		DoraAPIKey:  "key",
		Params:      map[string]string{"k": "v"},
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, k := range []string{
		"mode", "order_book_id", "start", "end", "resolution",
		"dora_base_url", "dora_api_key", "params",
	} {
		if _, ok := m[k]; !ok {
			t.Errorf("missing JSON key %q in %s", k, string(b))
		}
	}
}

func TestOrderIntentMirrorTags(t *testing.T) {
	intent := host.OrderIntent{
		Side:               "buy",
		Quantity:           "100",
		Type:               "limit",
		Price:              "99.5",
		InverseLeverage:    "2",
		FromGlobalPosition: false,
	}
	b, err := json.Marshal(intent)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, k := range []string{
		"side", "quantity", "type", "price", "inverse_leverage", "from_global_position",
	} {
		if _, ok := m[k]; !ok {
			t.Errorf("missing JSON key %q in %s", k, string(b))
		}
	}
}

func TestTradeMirrorTags(t *testing.T) {
	tr := host.Trade{
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
	b, err := json.Marshal(tr)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, k := range []string{
		"transaction_id", "order_book_id", "order_id", "order_seq",
		"user_id", "asset_0", "price", "quantity_0",
		"side", "aggressor_indicator", "created_at",
	} {
		if _, ok := jsonKey(b, k); !ok {
			t.Errorf("missing JSON key %q in %s", k, string(b))
		}
	}
}

func TestPriceMirrorTags(t *testing.T) {
	p := host.Price{
		AssetID: "asset-1",
		Price:   "100.5",
		YTM:     "0.0523",
		Time:    "2026-01-01T00:00:00Z",
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, k := range []string{"asset_id", "price", "ytm", "time"} {
		if _, ok := jsonKey(b, k); !ok {
			t.Errorf("missing JSON key %q in %s", k, string(b))
		}
	}
}

func TestCandleBatchMirrorTags(t *testing.T) {
	cb := host.CandleBatch{
		Items:  []host.Candle{{OrderBookID: "OB-1", Close: "100"}},
		Done:   true,
		Cursor: "abc",
	}
	b, err := json.Marshal(cb)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, k := range []string{"items", "done", "cursor"} {
		if _, ok := jsonKey(b, k); !ok {
			t.Errorf("missing JSON key %q in %s", k, string(b))
		}
	}
}

func TestTradeBatchMirrorTags(t *testing.T) {
	tb := host.TradeBatch{
		Items:  []host.Trade{{TransactionID: "tx-1"}},
		Done:   true,
		Cursor: "def",
	}
	b, err := json.Marshal(tb)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, k := range []string{"items", "done", "cursor"} {
		if _, ok := jsonKey(b, k); !ok {
			t.Errorf("missing JSON key %q in %s", k, string(b))
		}
	}
}

func TestPriceBatchMirrorTags(t *testing.T) {
	pb := host.PriceBatch{
		Items:  []host.Price{{AssetID: "asset-1"}},
		Done:   true,
		Cursor: "ghi",
	}
	b, err := json.Marshal(pb)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, k := range []string{"items", "done", "cursor"} {
		if _, ok := jsonKey(b, k); !ok {
			t.Errorf("missing JSON key %q in %s", k, string(b))
		}
	}
}

func TestEventEnvelopeRawData(t *testing.T) {
	// Data is json.RawMessage so the event's raw JSON object survives
	// a round-trip unchanged (a []byte field would base64-encode).
	ev := host.EventEnvelope{
		Type: "candle",
		Data: json.RawMessage(`{"order_book_id":"OB-1","close":"100"}`),
	}
	b, err := json.Marshal(ev)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got host.EventEnvelope
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Type != "candle" {
		t.Errorf("Type = %q, want %q", got.Type, "candle")
	}
	want := `{"order_book_id":"OB-1","close":"100"}`
	if string(got.Data) != want {
		t.Errorf("Data = %q, want %q (RawMessage must not base64-encode)", string(got.Data), want)
	}
}

// jsonKey unmarshals b into a generic map and reports whether key is present.
func jsonKey(b []byte, key string) (any, bool) {
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}
	v, ok := m[key]
	return v, ok
}
