package dorastrategy

import (
	"errors"
	"fmt"
)

// ErrModeNotImplemented is returned by Run when the requested mode is
// not yet implemented in this framework version. Plan 2 wires the
// wazero host runtime; Plan 3 fills in the backtest and live paths.
var ErrModeNotImplemented = errors.New("dorastrategy: mode not implemented")

// Run is the entrypoint a generated strategy's main() calls. It reads
// the config from the host via host.GetConfig(), dispatches on Mode,
// and returns the terminal error.
//
// Validate mode: calls s.Init(cfg) and returns. The config comes
// from the host's host_get_config import.
//
// Backtest mode: enters the candle loop (see runBacktest).
//
// Live mode: not yet implemented.
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
