package dorastrategy

// runBacktest is the ModeBacktest loop driver: the shared v3
// sequence (Init -> OnPreamble -> tagged event loop) with fills
// recorded for the backtest summary. Event ordering (including the
// warmup window) is decided server-side by host_next_event.
func runBacktest(s Strategy, cfg Config) error {
	return runEvents(s, cfg, true)
}
