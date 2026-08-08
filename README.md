# dora-strategy-wasm

Framework module for generated WASM strategy plugins. Sibling
sub-module of the [dora-agent](https://github.com/dora-network/dora-agent)
repo. Plan 1 ships the framework types, the host-call signatures,
the escape hatches, the manifest types, the TinyGo import
allowlist, and a build-time smoke. Plan 2 wires the server-side
wazero runtime to the host-call surface. Plan 3 fills in the
backtest and live modes. Plan 4 changes the LLM toolchain and
the validator to emit and accept the WASM path.

## Module path

```
github.com/dora-network/dora-strategy-wasm
```

The LLM-generated strategy imports:

- `github.com/dora-network/dora-strategy-wasm/dorastrategy` — framework types
- `github.com/dora-network/dora-strategy-wasm/dorastrategy/host` — host-call surface
- `github.com/dora-network/dora-strategy-wasm/dorastrategy/escape` — TinyGo escape hatches
- `github.com/dora-network/dora-strategy-wasm/manifest` — manifest struct (informational; the host reads the JSON directly)

## Build and test

```bash
go build ./...
go test ./...
```

The TinyGo build smoke compiles `example/strategy.go` to a valid
WASM module. CI skips the smoke if `tinygo` is not on PATH:

```bash
make tinygo-build
```

## What a generated strategy looks like

```go
package main

import (
    "github.com/dora-network/dora-strategy-wasm/dorastrategy"
    "github.com/dora-network/dora-strategy-wasm/dorastrategy/host"
)

type MyStrategy struct{}

func (s *MyStrategy) Init(cfg dorastrategy.Config) error {
    host.Log(host.LevelInfo, "init", "order_book", cfg.OrderBookID)
    return nil
}

func (s *MyStrategy) OnCandle(c dorastrategy.Candle) ([]dorastrategy.OrderIntent, error) {
    stop, _ := host.GetParamInt("stop_loss_bps")
    if c.Close > c.Open*1.005 {
        _ = stop
        return []dorastrategy.OrderIntent{{Side: "buy", Quantity: 100, Type: "market"}}, nil
    }
    return nil, nil
}

func main() {
    if err := dorastrategy.Run(&MyStrategy{}, dorastrategy.Config{Mode: dorastrategy.ModeValidate}); err != nil {
        panic(err)
    }
}
```

## Status

Plan 1. The server-side wiring (wazero runtime, read broker,
write broker, safety kernel) is in Plan 2+. The current
go-docker pipeline keeps working unchanged.
