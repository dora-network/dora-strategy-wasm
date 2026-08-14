// Package host is the host-call surface a generated WASM plugin uses.
// Each function here is backed by a wazero import that the server-side
// runtime wires at module instantiation time.
package host

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"unsafe"
)

const bufSize = 8192

var stackBuf [bufSize]byte //nolint:gochecknoglobals // stack buffer shared by all host calls in this package

// sentinelErr is the return value indicating a host-side error.
const sentinelErr int32 = -1

// Level is the severity of a log line.
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// OrderIntent is the host-side mirror of dorastrategy.OrderIntent. We
// re-declare it here to avoid an import cycle between the framework
// types and the wazero import surface. The fields mirror Dora's
// CreateOrderRequest (see dora-client-go and the dora-api OpenAPI
// spec); numeric fields are decimal strings to round-trip the
// wire format without precision loss.
type OrderIntent struct {
	Side               string `json:"side"`
	Quantity           string `json:"quantity"`
	Type               string `json:"type"`
	Price              string `json:"price"`
	InverseLeverage    string `json:"inverse_leverage"`
	FromGlobalPosition bool   `json:"from_global_position"`
}

// Candle is the host-side mirror of dorastrategy.Candle.
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

// Fill is the host-side mirror of dorastrategy.Fill.
type Fill struct {
	OrderID   string `json:"order_id"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Simulated bool   `json:"simulated"`
}

// Config is the host-side mirror of dorastrategy.Config.
type Config struct {
	Mode        string            `json:"mode"`
	OrderBookID string            `json:"order_book_id"`
	Start       time.Time         `json:"start"`
	End         time.Time         `json:"end"`
	Resolution  string            `json:"resolution"`
	DoraBaseURL string            `json:"dora_base_url"`
	DoraAPIKey  string            `json:"dora_api_key"`
	Params      map[string]string `json:"params"`
}

// Trade is the host-side mirror of dorastrategy.Trade. JSON tags
// match the wsplex /trades schema (see asyncapi.yaml).
type Trade struct {
	TransactionID      string `json:"transaction_id"`
	OrderBookID        string `json:"order_book_id"`
	OrderID            string `json:"order_id"`
	OrderSeq           int64  `json:"order_seq"`
	UserID             string `json:"user_id"`
	Asset0             string `json:"asset_0"` //nolint:tagliatelle // wsplex asyncapi wire name
	Price              string `json:"price"`
	Quantity0          string `json:"quantity_0"` //nolint:tagliatelle // wsplex asyncapi wire name
	Side               string `json:"side"`
	AggressorIndicator bool   `json:"aggressor_indicator"`
	CreatedAt          string `json:"created_at"`
}

// Price is the host-side mirror of dorastrategy.Price. JSON tags
// match the wsplex /prices schema (see asyncapi.yaml).
type Price struct {
	AssetID string `json:"asset_id"`
	Price   string `json:"price"`
	YTM     string `json:"ytm"`
	Time    string `json:"time"`
}

// CandleBatch is the host-side mirror of dorastrategy.CandleBatch.
type CandleBatch struct {
	Items  []Candle `json:"items"`
	Done   bool     `json:"done"`
	Cursor string   `json:"cursor"`
}

// TradeBatch is the host-side mirror of dorastrategy.TradeBatch.
type TradeBatch struct {
	Items  []Trade `json:"items"`
	Done   bool    `json:"done"`
	Cursor string  `json:"cursor"`
}

// PriceBatch is the host-side mirror of dorastrategy.PriceBatch.
type PriceBatch struct {
	Items  []Price `json:"items"`
	Done   bool    `json:"done"`
	Cursor string  `json:"cursor"`
}

// GetConfig returns the strategy's Config from the host.
func GetConfig() (Config, error) {
	n := wasmGetConfig(uint32(uintptr(unsafe.Pointer(&stackBuf[0]))), bufSize) //nolint:gosec // wazero ABI: pointer fits in u32
	if n == sentinelErr {
		return Config{}, readError()
	}
	if n <= 0 {
		return Config{}, errors.New("host: get_config returned no data")
	}
	var cfg Config
	if err := json.Unmarshal(stackBuf[:n], &cfg); err != nil {
		return Config{}, fmt.Errorf("host: decode config: %w", err)
	}
	return cfg, nil
}

// NextCandle returns the next candle, or (Candle{}, false) when done.
func NextCandle() (Candle, bool, error) {
	n := wasmNextCandle(uint32(uintptr(unsafe.Pointer(&stackBuf[0]))), bufSize) //nolint:gosec // wazero ABI: pointer fits in u32
	if n == sentinelErr {
		return Candle{}, false, readError()
	}
	if n == 0 {
		return Candle{}, false, nil
	}
	if n < 0 {
		return Candle{}, false, errors.New("host: next_candle returned invalid length")
	}
	var c Candle
	if err := json.Unmarshal(stackBuf[:n], &c); err != nil {
		return Candle{}, false, fmt.Errorf("host: decode candle: %w", err)
	}
	return c, true, nil
}

// NextLiveCandle returns the next live candle, blocking until one is available.
func NextLiveCandle() (Candle, bool, error) {
	n := wasmNextLiveCandle(uint32(uintptr(unsafe.Pointer(&stackBuf[0]))), bufSize) //nolint:gosec // wazero ABI: pointer fits in u32
	if n == sentinelErr {
		return Candle{}, false, readError()
	}
	if n == 0 {
		return Candle{}, false, nil
	}
	if n < 0 {
		return Candle{}, false, errors.New("host: next_live_candle returned invalid length")
	}
	var c Candle
	if err := json.Unmarshal(stackBuf[:n], &c); err != nil {
		return Candle{}, false, fmt.Errorf("host: decode live candle: %w", err)
	}
	return c, true, nil
}

// SubmitOrder sends an intent and returns the simulated fill.
func SubmitOrder(intent OrderIntent) (Fill, error) {
	intentJSON, err := json.Marshal(intent)
	if err != nil {
		return Fill{}, fmt.Errorf("host: encode intent: %w", err)
	}
	inLen := len(intentJSON)
	outOffset := bufSize / 2
	if inLen > outOffset {
		return Fill{}, errors.New("host: intent JSON exceeds half buffer")
	}
	copy(stackBuf[:inLen], intentJSON)
	n := wasmSubmitOrder(
		uint32(uintptr(unsafe.Pointer(&stackBuf[0]))), uint32(inLen), //nolint:gosec // wazero ABI
		uint32(uintptr(unsafe.Pointer(&stackBuf[outOffset]))), uint32(bufSize/2), //nolint:gosec,mnd // wazero ABI, half-buffer split
	)
	if n == sentinelErr {
		return Fill{}, readErrorOffset(outOffset)
	}
	if n <= 0 {
		return Fill{}, errors.New("host: submit_order returned no data")
	}
	var fill Fill
	if err := json.Unmarshal(stackBuf[outOffset:outOffset+int(n)], &fill); err != nil {
		return Fill{}, fmt.Errorf("host: decode fill: %w", err)
	}
	return fill, nil
}

// RecordFill records a fill for the summary.
func RecordFill(fill Fill) {
	fillJSON, err := json.Marshal(fill)
	if err != nil {
		return
	}
	copy(stackBuf[:len(fillJSON)], fillJSON)
	wasmRecordFill(uint32(uintptr(unsafe.Pointer(&stackBuf[0]))), uint32(len(fillJSON))) //nolint:gosec // wazero ABI
}

// Log emits a structured log line.
func Log(level Level, msg string) {
	var lvl uint32
	switch level {
	case LevelDebug:
		lvl = 0
	case LevelInfo:
		lvl = 1
	case LevelWarn:
		lvl = 2
	case LevelError:
		lvl = 3
	default:
		lvl = 1
	}
	b := []byte(msg)
	if len(b) == 0 {
		return
	}
	wasmLog(lvl, uint32(uintptr(unsafe.Pointer(&b[0]))), uint32(len(b))) //nolint:gosec // wazero ABI
}

// BacktestError signals a fatal error to the host.
func BacktestError(msg string) {
	b := []byte(msg)
	if len(b) == 0 {
		return
	}
	wasmBacktestError(uint32(uintptr(unsafe.Pointer(&b[0]))), uint32(len(b))) //nolint:gosec // wazero ABI
}

// readError reads the error JSON from the stack buffer. The host
// writes a 4-byte length prefix followed by the JSON when returning
// sentinelErr.
func readError() error {
	return readErrorOffset(0)
}

func readErrorOffset(offset int) error {
	if offset+4 > len(stackBuf) {
		return errors.New("host: error buffer too small")
	}
	errLen := int(binary.LittleEndian.Uint32(stackBuf[offset : offset+4]))
	if errLen <= 0 || errLen > len(stackBuf)-offset-4 {
		return errors.New("host: invalid error length")
	}
	var e struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(stackBuf[offset+4:offset+4+errLen], &e); err != nil {
		return errors.New("host: failed to decode error")
	}
	return errors.New(e.Error)
}

// fetchReq is the JSON envelope the host unmarshals from the
// plugin's input buffer for the three fetch imports.
type fetchReq struct {
	Start      string `json:"start,omitempty"`
	End        string `json:"end,omitempty"`
	Resolution string `json:"resolution,omitempty"`
	AssetID    string `json:"asset_id,omitempty"`
	BatchSize  int    `json:"batch_size"`
	Cursor     string `json:"cursor,omitempty"`
}

// FetchCandles is the host-side wrapper for the host_fetch_candles
// host call. The plugin's PreambleContext calls it; the real
// implementation reads from the history store. Tests override the
// package-level fetchCandlesFn.
func FetchCandles(req fetchReq) (CandleBatch, error) {
	return fetchCandlesFn(req)
}

// FetchTrades is the trade equivalent.
func FetchTrades(req fetchReq) (TradeBatch, error) {
	return fetchTradesFn(req)
}

// FetchPrices is the price equivalent.
func FetchPrices(req fetchReq) (PriceBatch, error) {
	return fetchPricesFn(req)
}

// EventEnvelope is one tagged event the plugin sees.
type EventEnvelope struct {
	Type string          `json:"type"` // "candle" | "trade" | "price"
	Data json.RawMessage `json:"data"` // raw JSON of the event
}

// NextEvent returns one tagged envelope (candle | trade |
// price) plus a `done` flag. The plugin's loop is a plain
// `for { ev, ok := host.NextEvent(); if !ok { break }; ...
// }`. The live host function does a `select` across three
// channels; the backtest host function walks the time-ordered
// union of the pre-fetched slices.
func NextEvent() (EventEnvelope, bool, error) {
	return nextEventFn()
}

// Test seams. In wasm builds these point to the real host
// imports; in tests they are replaced.
//
//nolint:gochecknoglobals // test seams; overridden in *_test.go
var (
	fetchCandlesFn = func(req fetchReq) (CandleBatch, error) {
		return CandleBatch{Done: true}, nil
	}
	fetchTradesFn = func(req fetchReq) (TradeBatch, error) {
		return TradeBatch{Done: true}, nil
	}
	fetchPricesFn = func(req fetchReq) (PriceBatch, error) {
		return PriceBatch{Done: true}, nil
	}
	nextEventFn = func() (EventEnvelope, bool, error) {
		return EventEnvelope{}, false, nil
	}
)
