package dorastrategy

import "errors"

// ErrModeNotImplemented is returned by Run when the requested mode is
// not yet implemented in this framework version. Plan 2 wires the
// wazero host runtime; Plan 3 fills in the backtest and live paths.
var ErrModeNotImplemented = errors.New("dorastrategy: mode not implemented")

// Run is the entrypoint a generated strategy's main() calls. It reads
// the mode from cfg, dispatches to the matching implementation, and
// returns the mode's terminal error (nil on success).
//
// Validate mode: calls s.Init(cfg) and returns its error. No network,
// no data. Network-free is a hard contract so the smoke test can run
// under --network=none.
//
// Backtest and Live modes: return ErrModeNotImplemented until Plan 3
// fills them in.
func Run(s Strategy, cfg Config) error {
	switch cfg.Mode {
	case ModeValidate:
		return s.Init(cfg)
	case ModeBacktest, ModeLive:
		return ErrModeNotImplemented
	default:
		return ErrModeNotImplemented
	}
}
