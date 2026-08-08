package dorastrategy_test

import (
	"errors"
	"testing"

	"github.com/dora-network/dora-strategy-wasm/dorastrategy"
)

func TestRun_ValidateMode_CallsInitAndReturnsNil(t *testing.T) {
	s := &recordingStrategy{}
	cfg := dorastrategy.Config{Mode: dorastrategy.ModeValidate}

	if err := dorastrategy.Run(s, cfg); err != nil {
		t.Fatalf("Run returned err: %v", err)
	}
	if !s.initCalled {
		t.Errorf("Init was not called")
	}
}

func TestRun_ValidateMode_PropagatesInitError(t *testing.T) {
	wantErr := errors.New("init failed")
	s := &recordingStrategy{initErr: wantErr}
	cfg := dorastrategy.Config{Mode: dorastrategy.ModeValidate}

	if err := dorastrategy.Run(s, cfg); !errors.Is(err, wantErr) {
		t.Fatalf("Run returned err=%v, want %v", err, wantErr)
	}
}

func TestRun_BacktestMode_NotImplemented(t *testing.T) {
	s := &recordingStrategy{}
	cfg := dorastrategy.Config{Mode: dorastrategy.ModeBacktest}

	err := dorastrategy.Run(s, cfg)
	if err == nil {
		t.Fatal("Run should return ErrModeNotImplemented for backtest")
	}
	if !errors.Is(err, dorastrategy.ErrModeNotImplemented) {
		t.Fatalf("Run returned err=%v, want ErrModeNotImplemented", err)
	}
}

func TestRun_LiveMode_NotImplemented(t *testing.T) {
	s := &recordingStrategy{}
	cfg := dorastrategy.Config{Mode: dorastrategy.ModeLive}

	if err := dorastrategy.Run(s, cfg); !errors.Is(err, dorastrategy.ErrModeNotImplemented) {
		t.Fatalf("Run returned err=%v, want ErrModeNotImplemented", err)
	}
}

type recordingStrategy struct {
	initCalled bool
	initErr    error
}

func (s *recordingStrategy) Init(dorastrategy.Config) error {
	s.initCalled = true
	return s.initErr
}

func (s *recordingStrategy) OnCandle(dorastrategy.Candle) ([]dorastrategy.OrderIntent, error) {
	return nil, nil
}
