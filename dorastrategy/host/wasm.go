//go:build wasm

package host

// These declarations are compiled by TinyGo as WASM imports the
// server (wazero) provides. They are excluded from regular `go build`
// via the `wasm` build constraint so the sub-module's unit tests
// compile without a WASM target.

//go:wasmimport env host_get_config
func wasmGetConfig(bufPtr uint32, bufLen uint32) int32

//go:wasmimport env host_next_candle
func wasmNextCandle(bufPtr uint32, bufLen uint32) int32

//go:wasmimport env host_next_live_candle
func wasmNextLiveCandle(bufPtr uint32, bufLen uint32) int32

//go:wasmimport env host_submit_order
func wasmSubmitOrder(inPtr uint32, inLen uint32, outPtr uint32, outLen uint32) int32

//go:wasmimport env host_record_fill
func wasmRecordFill(bufPtr uint32, bufLen uint32)

//go:wasmimport env host_log
func wasmLog(level uint32, bufPtr uint32, bufLen uint32)

//go:wasmimport env host_backtest_error
func wasmBacktestError(bufPtr uint32, bufLen uint32)

//go:wasmimport env host_fetch_candles
func wasmFetchCandles(inPtr uint32, inLen uint32, outPtr uint32, outLen uint32) int32

//go:wasmimport env host_fetch_trades
func wasmFetchTrades(inPtr uint32, inLen uint32, outPtr uint32, outLen uint32) int32

//go:wasmimport env host_fetch_prices
func wasmFetchPrices(inPtr uint32, inLen uint32, outPtr uint32, outLen uint32) int32

//go:wasmimport env host_next_event
func wasmNextEvent(outPtr uint32, outLen uint32) int32
