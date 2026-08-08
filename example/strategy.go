//go:build ignore

// Example strategy for the TinyGo build smoke. Build with:
//
//	tinygo build -target=wasi -buildmode=c-shared -o strategy.wasm strategy.go
//
// The build target is wired in the Makefile (Task 8). This file is
// excluded from the regular `go build` via the //go:build ignore
// tag so the parent module's `go build ./...` does not see it as a
// package.
//
// The strategy is intentionally trivial: it implements the
// dorastrategy.Strategy interface and returns nil from Init /
// OnCandle. Its only purpose is to prove the framework types
// compile against TinyGo and produce a valid .wasm module.
package main

import (
	"github.com/dora-network/dora-strategy-wasm/dorastrategy"
	"github.com/dora-network/dora-strategy-wasm/dorastrategy/host"
)

type noopStrategy struct{}

func (noopStrategy) Init(dorastrategy.Config) error { return nil }
func (noopStrategy) OnCandle(dorastrategy.Candle) ([]dorastrategy.OrderIntent, error) {
	_ = host.Log(host.LevelInfo, "noop on candle")
	return nil, nil
}

func main() {
	if err := dorastrategy.Run(noopStrategy{}, dorastrategy.Config{Mode: dorastrategy.ModeValidate}); err != nil {
		// framework's Run only returns nil for the validate path
		// (Init returns nil above). Panic is the right outcome if
		// the framework is misconfigured: the binary is a smoke
		// test, not a long-running process.
		panic(err)
	}
}
