package dorastrategy

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

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
	origNext := nextEventFn
	origSubmit := submitOrderFn
	origRecord := recordFillFn
	getConfigFn = func() (host.Config, error) {
		return host.Config{Mode: "backtest", OrderBookID: "OB-1"}, nil
	}
	calls := 0
	nextEventFn = func() (host.EventEnvelope, bool, error) {
		calls++
		if calls == 1 {
			return candleEvent(t, host.Candle{OrderBookID: "OB-1", Close: "100"}), true, nil
		}
		return host.EventEnvelope{}, false, nil
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
		nextEventFn = origNext
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
	origNext := nextEventFn
	origErr := backtestErrorFn
	getConfigFn = func() (host.Config, error) {
		return host.Config{Mode: "backtest"}, nil
	}
	nextEventFn = func() (host.EventEnvelope, bool, error) {
		return candleEvent(t, host.Candle{Close: "100"}), true, nil
	}
	var errMsg string
	backtestErrorFn = func(msg string) { errMsg = msg }
	defer func() {
		getConfigFn = origGetConfig
		nextEventFn = origNext
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

// candleEvent wraps a host candle in a tagged event envelope.
func candleEvent(t *testing.T, c host.Candle) host.EventEnvelope {
	t.Helper()
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal candle: %v", err)
	}
	return host.EventEnvelope{Type: "candle", Data: data}
}

type recordingStrategy struct {
	StrategyBase
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

// TestRun_DispatchesOnPreamble proves the v3 sequence end to end:
// Init -> OnPreamble (whose FetchCandles paginates through the host
// seam) -> OnCandle fed by the event loop.
func TestRun_DispatchesOnPreamble(t *testing.T) {
	origGetConfig := getConfigFn
	origFetch := fetchCandlesFn
	origNext := nextEventFn
	origSubmit := submitOrderFn
	origRecord := recordFillFn
	getConfigFn = func() (host.Config, error) {
		return host.Config{Mode: "backtest", OrderBookID: "OB-1"}, nil
	}
	var reqs []host.FetchReq
	fetchCandlesFn = func(req host.FetchReq) (host.CandleBatch, error) {
		reqs = append(reqs, req)
		if len(reqs) == 1 {
			return host.CandleBatch{
				Items:  []host.Candle{{OrderBookID: "OB-1", Close: "99"}},
				Cursor: "cur-1",
			}, nil
		}
		return host.CandleBatch{Done: true}, nil
	}
	nextEventDone := false
	nextEventFn = func() (host.EventEnvelope, bool, error) {
		if nextEventDone {
			return host.EventEnvelope{}, false, nil
		}
		nextEventDone = true
		return candleEvent(t, host.Candle{OrderBookID: "OB-1", Close: "100"}), true, nil
	}
	submitOrderFn = func(host.OrderIntent) (host.Fill, error) {
		return host.Fill{}, nil
	}
	recordFillFn = func(host.Fill) {}
	defer func() {
		getConfigFn = origGetConfig
		fetchCandlesFn = origFetch
		nextEventFn = origNext
		submitOrderFn = origSubmit
		recordFillFn = origRecord
	}()

	s := &preambleStrategy{}
	if err := Run(s); err != nil {
		t.Fatalf("Run returned err: %v", err)
	}
	if !s.initCalled {
		t.Errorf("Init was not called")
	}
	if !s.preambleCalled {
		t.Errorf("OnPreamble was not called")
	}
	if len(s.warmCandles) != 1 || s.warmCandles[0].Close != "99" {
		t.Errorf("warm candles: got %+v", s.warmCandles)
	}
	if len(reqs) != 2 {
		t.Fatalf("FetchCandles called %d times, want 2", len(reqs))
	}
	if reqs[0].Cursor != "" || reqs[1].Cursor != "cur-1" {
		t.Errorf("cursor threading: got %+v", reqs)
	}
	if reqs[0].Start == "" || reqs[0].End == "" || reqs[0].Resolution != "1m" {
		t.Errorf("fetch request fields: got %+v", reqs[0])
	}
	if !s.candleCalled {
		t.Errorf("OnCandle was not called after OnPreamble")
	}
}

// preambleStrategy records the v3 hook sequence and warms up via
// FetchCandles.
type preambleStrategy struct {
	StrategyBase
	initCalled     bool
	preambleCalled bool
	candleCalled   bool
	warmCandles    []Candle
}

func (s *preambleStrategy) Init(Config) error {
	s.initCalled = true
	return nil
}

func (s *preambleStrategy) OnPreamble(ctx context.Context, p PreambleContext) error {
	s.preambleCalled = true
	for {
		batch, err := p.FetchCandles(ctx,
			time.Now().Add(-7*24*time.Hour), time.Now(), "1m", 200)
		if err != nil {
			return err
		}
		s.warmCandles = append(s.warmCandles, batch.Items...)
		if batch.Done {
			return nil
		}
	}
}

func (s *preambleStrategy) OnCandle(Candle) ([]OrderIntent, error) {
	s.candleCalled = true
	return nil, nil
}
