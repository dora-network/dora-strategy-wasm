package dorastrategy

import (
	"github.com/dora-network/dora-strategy-wasm/dorastrategy/host"
)

// runBacktest is the ModeBacktest loop driver. The plugin's Run
// function calls this after reading the config from the host.
// The loop: pull candles from the host, feed to OnCandle, submit
// any returned intents, record fills. Exits when host_next_candle
// signals done (returns ok=false).
func runBacktest(s Strategy, cfg Config) error {
	if err := s.Init(cfg); err != nil {
		backtestErrorFn(err.Error())
		return err
	}
	for {
		hostCandle, ok, err := nextCandleFn()
		if err != nil {
			backtestErrorFn(err.Error())
			return err
		}
		if !ok {
			break
		}
		intents, err := s.OnCandle(Candle(hostCandle))
		if err != nil {
			backtestErrorFn(err.Error())
			return err
		}
		for _, intent := range intents {
			fill, err := submitOrderFn(host.OrderIntent(intent))
			if err != nil {
				backtestErrorFn(err.Error())
				return err
			}
			recordFillFn(fill)
		}
	}
	return nil
}

// Test seams. In wasm builds these point to the real host imports;
// in unit tests they are replaced with fakes so the loop logic can be
// exercised without a wazero runtime.
//
//nolint:gochecknoglobals // test seams; overridden in *_test.go
var (
	getConfigFn     = host.GetConfig
	nextCandleFn    = host.NextCandle
	submitOrderFn   = host.SubmitOrder
	recordFillFn    = host.RecordFill
	backtestErrorFn = host.BacktestError
)
