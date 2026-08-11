package dorastrategy

import (
	"github.com/dora-network/dora-strategy-wasm/dorastrategy/host"
)

// runLive is the ModeLive loop driver. Unlike backtest, candles
// arrive in real time from the host; the loop blocks on
// nextLiveCandleFn until a candle is published. Fills are not
// recorded here — in live mode fills come from the Dora REST API
// and are observed via the orders-update stream, not the candle
// loop. Orders are still submitted through submitOrderFn.
func runLive(s Strategy, cfg Config) error {
	if err := s.Init(cfg); err != nil {
		backtestErrorFn(err.Error())
		return err
	}
	for {
		hostCandle, ok, err := nextLiveCandleFn()
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
			if _, err := submitOrderFn(host.OrderIntent(intent)); err != nil {
				backtestErrorFn(err.Error())
				return err
			}
		}
	}
	return nil
}
