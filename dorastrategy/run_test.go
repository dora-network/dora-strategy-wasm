package dorastrategy

import (
	"errors"
	"testing"

	"github.com/dora-network/dora-strategy-wasm/dorastrategy/host"
)

func TestRun_ValidateMode_CallsInitAndReturnsNil(t *testing.T) {
	orig := getConfigFn
	getConfigFn = func() (host.Config, error) {
		return host.Config{Mode: "validate"}, nil
	}
	defer func() { getConfigFn = orig }()

	s := &recordingStrategy{}
	if err := Run(s); err != nil {
		t.Fatalf("Run returned err: %v", err)
	}
	if !s.initCalled {
		t.Errorf("Init was not called")
	}
}

func TestRun_ValidateMode_PropagatesInitError(t *testing.T) {
	orig := getConfigFn
	wantErr := errors.New("init failed")
	getConfigFn = func() (host.Config, error) {
		return host.Config{Mode: "validate"}, nil
	}
	defer func() { getConfigFn = orig }()

	s := &recordingStrategy{initErr: wantErr}
	if err := Run(s); !errors.Is(err, wantErr) {
		t.Fatalf("Run returned err=%v, want %v", err, wantErr)
	}
}

func TestRun_GetConfigError(t *testing.T) {
	orig := getConfigFn
	wantErr := errors.New("config unavailable")
	getConfigFn = func() (host.Config, error) {
		return host.Config{}, wantErr
	}
	defer func() { getConfigFn = orig }()

	err := Run(&recordingStrategy{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("Run returned err=%v, want %v", err, wantErr)
	}
}

func TestRun_BacktestMode_RunsLoop(t *testing.T) {
	origGetConfig := getConfigFn
	origNext := nextCandleFn
	origSubmit := submitOrderFn
	origRecord := recordFillFn
	getConfigFn = func() (host.Config, error) {
		return host.Config{Mode: "backtest", OrderBookID: "OB-1"}, nil
	}
	calls := 0
	nextCandleFn = func() (host.Candle, bool, error) {
		calls++
		if calls == 1 {
			return host.Candle{OrderBookID: "OB-1", Close: "100"}, true, nil
		}
		return host.Candle{}, false, nil
	}
	var submitted host.OrderIntent
	submitOrderFn = func(intent host.OrderIntent) (host.Fill, error) {
		submitted = intent
		return host.Fill{OrderID: "sim-1", Price: "100", Quantity: "10", Simulated: true}, nil
	}
	var recorded []host.Fill
	recordFillFn = func(fill host.Fill) {
		recorded = append(recorded, fill)
	}
	defer func() {
		getConfigFn = origGetConfig
		nextCandleFn = origNext
		submitOrderFn = origSubmit
		recordFillFn = origRecord
	}()

	s := &recordingStrategy{onCandleIntents: []OrderIntent{{Side: "buy", Quantity: "10", Type: "market"}}}
	if err := Run(s); err != nil {
		t.Fatalf("Run returned err: %v", err)
	}
	if !s.initCalled {
		t.Errorf("Init was not called")
	}
	if !s.candleCalled {
		t.Errorf("OnCandle was not called")
	}
	if submitted.Quantity != "10" {
		t.Errorf("submit order intent: got %+v", submitted)
	}
	if len(recorded) != 1 || recorded[0].OrderID != "sim-1" {
		t.Errorf("recorded fills: got %+v", recorded)
	}
}

func TestRun_BacktestMode_PropagatesOnCandleError(t *testing.T) {
	origGetConfig := getConfigFn
	origNext := nextCandleFn
	origErr := backtestErrorFn
	getConfigFn = func() (host.Config, error) {
		return host.Config{Mode: "backtest"}, nil
	}
	nextCandleFn = func() (host.Candle, bool, error) {
		return host.Candle{Close: "100"}, true, nil
	}
	var errMsg string
	backtestErrorFn = func(msg string) { errMsg = msg }
	defer func() {
		getConfigFn = origGetConfig
		nextCandleFn = origNext
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

type recordingStrategy struct {
	initCalled      bool
	initErr         error
	candleCalled    bool
	onCandleIntents []OrderIntent
	onCandleErr     error
}

func (s *recordingStrategy) Init(Config) error {
	s.initCalled = true
	return s.initErr
}

func (s *recordingStrategy) OnCandle(Candle) ([]OrderIntent, error) {
	s.candleCalled = true
	return s.onCandleIntents, s.onCandleErr
}
