package dorastrategy

// runLive is the ModeLive loop driver: the same v3 sequence as
// backtest, minus fill recording — in live mode fills come from the
// Dora REST API and are observed via the orders-update stream, not
// the event loop. Orders are still submitted through the host.
func runLive(s Strategy, cfg Config) error {
	return runEvents(s, cfg, false)
}
