package dorastrategy

import (
	"context"
	"time"

	"github.com/dora-network/dora-strategy-wasm/dorastrategy/host"
)

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

// Trade is the canonical trade record. JSON tags match the wsplex
// /trades schema (see asyncapi.yaml: transaction_id, asset_0,
// quantity_0, aggressor_indicator, order_seq, user_id,
// order_book_id). The live host function passes the wsplex
// notification through unchanged; the backtest / preamble
// FetchTrades SELECT aliases trades_history columns onto the
// same names (§4.4 of the spec).
type Trade struct {
	TransactionID      string `json:"transaction_id"`
	OrderBookID        string `json:"order_book_id"`
	OrderID            string `json:"order_id"`
	OrderSeq           int64  `json:"order_seq"`
	UserID             string `json:"user_id"`
	Asset0             string `json:"asset_0"` //nolint:tagliatelle // wsplex asyncapi wire name
	Price              string `json:"price"`
	Quantity0          string `json:"quantity_0"`          //nolint:tagliatelle // wsplex asyncapi wire name
	Side               string `json:"side"`                // "BUY" | "SELL"
	AggressorIndicator bool   `json:"aggressor_indicator"` // true = taker
	CreatedAt          string `json:"created_at"`
}

// Price is the canonical price record. JSON tags match the wsplex
// /prices schema (see asyncapi.yaml: asset_id, price, ytm, time).
type Price struct {
	AssetID string `json:"asset_id"`
	Price   string `json:"price"`
	YTM     string `json:"ytm"`
	Time    string `json:"time"`
}

// FrameworkVersion is the tagged version of the dorastrategy
// framework. The agent (dora-agent) does not know the framework
// version until the framework is tagged, so this defaults to
// "dev" and is overwritten at validate-time by the agent's
// tinygo build via -ldflags:
//
//	tinygo build -ldflags \
//	  "-X github.com/dora-network/dora-strategy-wasm/dorastrategy.FrameworkVersion=<tag>"
//
// The agent's WasmFrameworkVersion constant (in
// internal/llm/prompts) is bumped per release to match the
// framework's tagged release.
var FrameworkVersion = "dev" //nolint:gochecknoglobals // build-time injected via -ldflags, not state.

type PreambleContext interface {
	FetchCandles(ctx context.Context, start, end time.Time, resolution string, batchSize int) (CandleBatch, error)
	FetchTrades(ctx context.Context, start, end time.Time, batchSize int) (TradeBatch, error)
	FetchPrices(ctx context.Context, start, end time.Time, batchSize int) (PriceBatch, error)
}

// CandleBatch is one batch of historic candles returned by
// FetchCandles. The plugin passes Cursor back to the next call
// until Done is true.
type CandleBatch struct {
	Items  []Candle
	Done   bool
	Cursor string
}

// TradeBatch is one batch of historic trades returned by
// FetchTrades. Trades are not bucketed.
type TradeBatch struct {
	Items  []Trade
	Done   bool
	Cursor string
}

// PriceBatch is one batch of historic prices returned by
// FetchPrices. Prices are per-event records, not bucketed.
type PriceBatch struct {
	Items  []Price
	Done   bool
	Cursor string
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
// Mode-independent: Init runs once at startup; OnPreamble runs
// once at startup (network-bound, populates state); OnCandle /
// OnTrade / OnPrice run per-event in both backtest and live modes
// (the order is arrival time, including the warmup window for
// backtest). All hooks except OnPreamble may return OrderIntents
// to place orders. OnPreamble's return type is just `error` so
// the type system enforces that the preamble never submits
// orders (§4.2 of the spec).
type Strategy interface {
	Init(cfg Config) error
	OnPreamble(ctx context.Context, p PreambleContext) error
	OnCandle(c Candle) ([]OrderIntent, error)
	OnTrade(t Trade) ([]OrderIntent, error)
	OnPrice(p Price) ([]OrderIntent, error)
}

// StrategyBase is a convenience embed for plugins that don't
// care about OnPreamble, OnTrade, or OnPrice. It provides no-op
// defaults so an LLM-generated strategy can override only the
// hooks it needs. A plugin that wants trades/prices MUST
// override the corresponding method to return early
// (`return nil, nil`).
type StrategyBase struct{}

func (StrategyBase) Init(Config) error { return nil }
func (StrategyBase) OnPreamble(context.Context, PreambleContext) error {
	return nil
}
func (StrategyBase) OnCandle(Candle) ([]OrderIntent, error) { return nil, nil }
func (StrategyBase) OnTrade(Trade) ([]OrderIntent, error)   { return nil, nil }
func (StrategyBase) OnPrice(Price) ([]OrderIntent, error)   { return nil, nil }

// preambleCtx is the concrete PreambleContext the framework hands to
// OnPreamble. It wraps the host fetch imports and threads the
// pagination cursor across sequential calls.
//
// ponytail: cursors live on the struct because the committed v3
// PreambleContext interface has no cursor parameter; OnPreamble is
// single-threaded and strictly sequential, so this holds. Revisit
// only if the interface ever gains a cursor param.
type preambleCtx struct {
	candleCursor string
	tradeCursor  string
	priceCursor  string
}

func (p *preambleCtx) FetchCandles(ctx context.Context, start, end time.Time, resolution string, batchSize int) (CandleBatch, error) {
	if err := ctx.Err(); err != nil {
		return CandleBatch{}, err
	}
	b, err := fetchCandlesFn(host.FetchReq{
		Start:      start.Format(time.RFC3339),
		End:        end.Format(time.RFC3339),
		Resolution: resolution,
		BatchSize:  batchSize,
		Cursor:     p.candleCursor,
	})
	if err != nil {
		return CandleBatch{}, err
	}
	p.candleCursor = b.Cursor
	out := CandleBatch{Done: b.Done, Cursor: b.Cursor}
	for _, c := range b.Items {
		out.Items = append(out.Items, Candle(c))
	}
	return out, nil
}

func (p *preambleCtx) FetchTrades(ctx context.Context, start, end time.Time, batchSize int) (TradeBatch, error) {
	if err := ctx.Err(); err != nil {
		return TradeBatch{}, err
	}
	b, err := fetchTradesFn(host.FetchReq{
		Start:     start.Format(time.RFC3339),
		End:       end.Format(time.RFC3339),
		BatchSize: batchSize,
		Cursor:    p.tradeCursor,
	})
	if err != nil {
		return TradeBatch{}, err
	}
	p.tradeCursor = b.Cursor
	out := TradeBatch{Done: b.Done, Cursor: b.Cursor}
	for _, tr := range b.Items {
		out.Items = append(out.Items, Trade(tr))
	}
	return out, nil
}

func (p *preambleCtx) FetchPrices(ctx context.Context, start, end time.Time, batchSize int) (PriceBatch, error) {
	if err := ctx.Err(); err != nil {
		return PriceBatch{}, err
	}
	b, err := fetchPricesFn(host.FetchReq{
		Start:     start.Format(time.RFC3339),
		End:       end.Format(time.RFC3339),
		BatchSize: batchSize,
		Cursor:    p.priceCursor,
	})
	if err != nil {
		return PriceBatch{}, err
	}
	p.priceCursor = b.Cursor
	out := PriceBatch{Done: b.Done, Cursor: b.Cursor}
	for _, pr := range b.Items {
		out.Items = append(out.Items, Price(pr))
	}
	return out, nil
}
