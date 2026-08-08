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

// Candle is a mode-neutral OHLCV bar. Same shape for backtest replay
// and live-aggregated ticks.
type Candle struct {
	Timestamp                      time.Time
	Open, High, Low, Close, Volume float64
}

// OrderIntent is the strategy's desire to trade. Mode-neutral; the
// framework translates it into a real or simulated order.
type OrderIntent struct {
	Side     string // "buy" | "sell"
	Quantity float64
	Type     string  // "market" | "limit"
	Price    float64 // limit orders only
}

// Fill is the result of submitting an intent (real or simulated).
type Fill struct {
	OrderID         string
	Price, Quantity float64
	Simulated       bool
}

// Config is parsed by the framework from env/flags.
type Config struct {
	Mode        Mode
	OrderBookID string
	Start, End  time.Time // backtest window
	Resolution  string    // candle resolution
	DoraBaseURL string
	DoraAPIKey  string
	Params      map[string]string // strategy-specific
}

// Strategy is the decision core the generated code implements.
// Mode-independent: Init runs once at startup; OnCandle is invoked per
// candle (historic replay or live-aggregated). Init MUST be network-free
// so that validate mode works under --network=none.
type Strategy interface {
	Init(cfg Config) error
	OnCandle(c Candle) ([]OrderIntent, error)
}
