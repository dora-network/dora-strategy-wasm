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
