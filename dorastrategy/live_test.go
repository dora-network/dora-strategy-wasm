package dorastrategy

import (
	"errors"
	"testing"

	"github.com/dora-network/dora-strategy-wasm/dorastrategy/host"
)

func TestRunLive_HappyPath(t *testing.T) {
	origGetConfig := getConfigFn
	origNext := nextLiveCandleFn
	origSubmit := submitOrderFn
	getConfigFn = func() (host.Config, error) {
		return host.Config{Mode: "live", OrderBookID: "OB-1"}, nil
	}
	calls := 0
	nextLiveCandleFn = func() (host.Candle, bool, error) {
		calls++
		switch calls {
		case 1:
			return host.Candle{OrderBookID: "OB-1", Close: "100"}, true, nil
		case 2:
			return host.Candle{OrderBookID: "OB-1", Close: "101"}, true, nil
		default:
			return host.Candle{}, false, nil
		}
	}
	var submitted []host.OrderIntent
	submitOrderFn = func(intent host.OrderIntent) (host.Fill, error) {
		submitted = append(submitted, intent)
		return host.Fill{OrderID: "live-1"}, nil
	}
	defer func() {
		getConfigFn = origGetConfig
		nextLiveCandleFn = origNext
		submitOrderFn = origSubmit
	}()

	s := &recordingStrategy{onCandleIntents: []OrderIntent{{Side: "buy", Quantity: "5", Type: "market"}}}
	if err := Run(s); err != nil {
		t.Fatalf("Run returned err: %v", err)
	}
	if !s.initCalled {
		t.Errorf("Init was not called")
	}
	if calls != 3 {
		t.Errorf("nextLiveCandle called %d times, want 3 (2 candles + 1 EOF)", calls)
	}
	if len(submitted) != 2 {
		t.Errorf("submitOrder called %d times, want 2", len(submitted))
	}
}

func TestRunLive_InitError(t *testing.T) {
	origGetConfig := getConfigFn
	origErr := backtestErrorFn
	getConfigFn = func() (host.Config, error) {
		return host.Config{Mode: "live"}, nil
	}
	var errMsg string
	backtestErrorFn = func(msg string) { errMsg = msg }
	defer func() {
		getConfigFn = origGetConfig
		backtestErrorFn = origErr
	}()

	wantErr := errors.New("init boom")
	s := &recordingStrategy{initErr: wantErr}
	if err := Run(s); !errors.Is(err, wantErr) {
		t.Fatalf("Run returned err=%v, want %v", err, wantErr)
	}
	if errMsg != wantErr.Error() {
		t.Errorf("BacktestError msg=%q want %q", errMsg, wantErr.Error())
	}
}

func TestRunLive_OnCandleError(t *testing.T) {
	origGetConfig := getConfigFn
	origNext := nextLiveCandleFn
	origErr := backtestErrorFn
	getConfigFn = func() (host.Config, error) {
		return host.Config{Mode: "live"}, nil
	}
	nextLiveCandleFn = func() (host.Candle, bool, error) {
		return host.Candle{Close: "100"}, true, nil
	}
	var errMsg string
	backtestErrorFn = func(msg string) { errMsg = msg }
	defer func() {
		getConfigFn = origGetConfig
		nextLiveCandleFn = origNext
		backtestErrorFn = origErr
	}()

	wantErr := errors.New("candle error")
	s := &recordingStrategy{onCandleErr: wantErr}
	if err := Run(s); !errors.Is(err, wantErr) {
		t.Fatalf("Run returned err=%v, want %v", err, wantErr)
	}
	if errMsg != wantErr.Error() {
		t.Errorf("BacktestError msg=%q want %q", errMsg, wantErr.Error())
	}
}
