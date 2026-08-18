package dorastrategy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dora-network/dora-strategy-wasm/dorastrategy/host"
)

// ErrModeNotImplemented is returned by Run when the requested mode is
// not yet implemented in this framework version.
var ErrModeNotImplemented = errors.New("dorastrategy: mode not implemented")

// Run is the entrypoint a generated strategy's main() calls. It reads
// the config from the host via host.GetConfig(), dispatches on Mode,
// and returns the terminal error.
//
// Validate mode: calls s.Init(cfg) and returns.
//
// Backtest and live mode: runEvents — Init, OnPreamble (warm-up
// fetches via PreambleContext), then the tagged event loop
// (OnCandle / OnTrade / OnPrice) fed by host.NextEvent.
func Run(s Strategy) error {
	hc, err := getConfigFn()
	if err != nil {
		return fmt.Errorf("dorastrategy: get config: %w", err)
	}
	cfg := Config{
		Mode:        Mode(hc.Mode),
		OrderBookID: hc.OrderBookID,
		Start:       hc.Start,
		End:         hc.End,
		Resolution:  hc.Resolution,
		DoraBaseURL: hc.DoraBaseURL,
		DoraAPIKey:  hc.DoraAPIKey,
		Params:      hc.Params,
	}
	switch cfg.Mode {
	case ModeValidate:
		return s.Init(cfg)
	case ModeBacktest:
		return runBacktest(s, cfg)
	case ModeLive:
		return runLive(s, cfg)
	default:
		return ErrModeNotImplemented
	}
}

// runEvents is the shared backtest/live driver: Init, OnPreamble,
// then the event loop. recordFills is true in backtest mode (fills
// feed the summary); in live mode fills arrive via the orders-update
// stream, so only submission happens here.
func runEvents(s Strategy, cfg Config, recordFills bool) error {
	if err := s.Init(cfg); err != nil {
		backtestErrorFn(err.Error())
		return err
	}
	if err := s.OnPreamble(context.Background(), &preambleCtx{}); err != nil {
		backtestErrorFn(err.Error())
		return err
	}
	for {
		ev, ok, err := nextEventFn()
		if err != nil {
			backtestErrorFn(err.Error())
			return err
		}
		if !ok {
			break
		}
		intents, err := dispatchEvent(s, ev)
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
			if recordFills {
				recordFillFn(fill)
			}
		}
	}
	return nil
}

// dispatchEvent decodes one tagged envelope and routes it to the
// matching Strategy hook.
func dispatchEvent(s Strategy, ev host.EventEnvelope) ([]OrderIntent, error) {
	switch ev.Type {
	case "candle":
		var hc host.Candle
		if err := json.Unmarshal(ev.Data, &hc); err != nil {
			return nil, fmt.Errorf("dorastrategy: decode candle event: %w", err)
		}
		return s.OnCandle(Candle(hc))
	case "trade":
		var ht host.Trade
		if err := json.Unmarshal(ev.Data, &ht); err != nil {
			return nil, fmt.Errorf("dorastrategy: decode trade event: %w", err)
		}
		return s.OnTrade(Trade(ht))
	case "price":
		var hp host.Price
		if err := json.Unmarshal(ev.Data, &hp); err != nil {
			return nil, fmt.Errorf("dorastrategy: decode price event: %w", err)
		}
		return s.OnPrice(Price(hp))
	default:
		return nil, fmt.Errorf("dorastrategy: unknown event type %q", ev.Type)
	}
}

// Test seams. In wasm builds these point to the real host imports;
// in unit tests they are replaced with fakes so the loop logic can be
// exercised without a wazero runtime.
//
//nolint:gochecknoglobals // test seams; overridden in *_test.go
var (
	getConfigFn     = host.GetConfig
	nextEventFn     = host.NextEvent
	fetchCandlesFn  = host.FetchCandles
	fetchTradesFn   = host.FetchTrades
	fetchPricesFn   = host.FetchPrices
	submitOrderFn   = host.SubmitOrder
	recordFillFn    = host.RecordFill
	backtestErrorFn = host.BacktestError
)
