//go:build !wasm

package host

// Stub implementations for non-wasm builds (unit tests). These panic
// if called; in tests the host functions are never invoked because
// Run returns ErrModeNotImplemented for non-validate modes.

func wasmGetConfig(bufPtr uint32, bufLen uint32) int32 {
	panic("host: wasm import called in non-wasm build")
}

func wasmNextCandle(bufPtr uint32, bufLen uint32) int32 {
	panic("host: wasm import called in non-wasm build")
}

func wasmNextLiveCandle(_, _ uint32) int32 { return 0 }

func wasmSubmitOrder(inPtr, inLen, outPtr, outLen uint32) int32 {
	panic("host: wasm import called in non-wasm build")
}

func wasmRecordFill(bufPtr uint32, bufLen uint32) {}

func wasmLog(level, bufPtr, bufLen uint32) {}

func wasmBacktestError(bufPtr uint32, bufLen uint32) {}

// Each pair corresponds to a //go:wasmimport declaration in wasm.go
// compiled only under -tags wasm; these stubs satisfy the !wasm build
// so the package compiles in unit tests.
func wasmFetchCandles(inPtr, inLen, outPtr, outLen uint32) int32 {
	return 0
}

func wasmFetchTrades(inPtr, inLen, outPtr, outLen uint32) int32 {
	return 0
}

func wasmFetchPrices(inPtr, inLen, outPtr, outLen uint32) int32 {
	return 0
}

func wasmNextEvent(outPtr, outLen uint32) int32 {
	return 0
}
