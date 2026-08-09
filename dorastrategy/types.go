package dorastrategy

import "time"

// Mode selects the runtime mode of a strategy instance. The framework
// dispatches to the matching implementation; today only ModeValidate is
// fully implemented (see run.go).
type Mode string

const (
	ModeValidate Mode = "validate" // init + exit 0 (smoke test, no network)
	ModeBacktest Mode = "backtest" // historic replay + simulated orders
	ModeLive     Mode = "live"     // real-time feeds + real orders
)

// Candle carries every field Dora's SDK returns for a candlestick
// bar. All numeric fields are decimal strings to preserve precision
// (matching the SDK wire format). StartTimestamp is RFC3339.
type Candle struct {
	OrderBookID    string `json:"order_book_id"`
	StartTimestamp string `json:"start_timestamp"`
	Open           string `json:"open"`
	High           string `json:"high"`
	Low            string `json:"low"`
	Close          string `json:"close"`
	OpenYtm        string `json:"open_ytm"`
	CloseYtm       string `json:"close_ytm"`
	HighYtm        string `json:"high_ytm"`
	LowYtm         string `json:"low_ytm"`
	Volume         string `json:"volume"`
}

// OrderIntent is the strategy's desire to trade. Mode-neutral; the
// framework translates it into a real or simulated order.
//
// Field shapes mirror Dora's CreateOrderRequest (see
// github.com/dora-network/dora-client-go/doraclient.CreateOrderRequest and
// the dora-api OpenAPI spec for the createOrder operation). Numeric
// fields are decimal strings, not float64, to round-trip through
// the wire format without precision loss.
type OrderIntent struct {
	Side               string `json:"side"`                 // "buy" | "sell"
	Quantity           string `json:"quantity"`             // decimal string, e.g. "100", "0.5"
	Type               string `json:"type"`                 // "market" | "limit"
	Price              string `json:"price"`                // required for limit, empty for market
	InverseLeverage    string `json:"inverse_leverage"`     // decimal string; empty => host default "1" (no leverage)
	FromGlobalPosition bool   `json:"from_global_position"` // false (default) => isolated margin; true => global cash account
}

// Fill is the result of submitting an intent (real or simulated).
type Fill struct {
	OrderID   string `json:"order_id"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Simulated bool   `json:"simulated"`
}

// Config is parsed by the framework from env/flags.
type Config struct {
	Mode        Mode              `json:"mode"`
	OrderBookID string            `json:"order_book_id"`
	Start       time.Time         `json:"start"` // backtest window
	End         time.Time         `json:"end"`
	Resolution  string            `json:"resolution"`
	DoraBaseURL string            `json:"dora_base_url"`
	DoraAPIKey  string            `json:"dora_api_key"`
	Params      map[string]string `json:"params"`
}

// Strategy is the decision core the generated code implements.
// Mode-independent: Init runs once at startup; OnCandle is invoked per
// candle (historic replay or live-aggregated). Init MUST be network-free
// so that validate mode works under --network=none.
type Strategy interface {
	Init(cfg Config) error
	OnCandle(c Candle) ([]OrderIntent, error)
}
