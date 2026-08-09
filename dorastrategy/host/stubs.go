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

func wasmSubmitOrder(inPtr, inLen, outPtr, outLen uint32) int32 {
	panic("host: wasm import called in non-wasm build")
}

func wasmRecordFill(bufPtr uint32, bufLen uint32) {}

func wasmLog(level, bufPtr, bufLen uint32) {}

func wasmBacktestError(bufPtr uint32, bufLen uint32) {}
