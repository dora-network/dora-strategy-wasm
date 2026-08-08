package host_test

import (
	"errors"
	"testing"

	"github.com/dora-network/dora-strategy-wasm/dorastrategy/host"
)

func TestLog_NotWired_ReturnsErr(t *testing.T) {
	err := host.Log(host.LevelInfo, "test")
	if !errors.Is(err, host.ErrNotWired) {
		t.Fatalf("Log returned err=%v, want ErrNotWired", err)
	}
}

func TestGetParamString_NotWired(t *testing.T) {
	_, err := host.GetParamString("k")
	if !errors.Is(err, host.ErrNotWired) {
		t.Fatalf("GetParamString returned err=%v, want ErrNotWired", err)
	}
}

func TestGetParamInt_NotWired(t *testing.T) {
	_, err := host.GetParamInt("k")
	if !errors.Is(err, host.ErrNotWired) {
		t.Fatalf("GetParamInt returned err=%v, want ErrNotWired", err)
	}
}

func TestGetParamFloat_NotWired(t *testing.T) {
	_, err := host.GetParamFloat("k")
	if !errors.Is(err, host.ErrNotWired) {
		t.Fatalf("GetParamFloat returned err=%v, want ErrNotWired", err)
	}
}

func TestGetParamBool_NotWired(t *testing.T) {
	_, err := host.GetParamBool("k")
	if !errors.Is(err, host.ErrNotWired) {
		t.Fatalf("GetParamBool returned err=%v, want ErrNotWired", err)
	}
}

func TestSubmitOrder_NotWired(t *testing.T) {
	_, err := host.SubmitOrder(host.OrderIntent{Side: "buy", Quantity: 1, Type: "market"})
	if !errors.Is(err, host.ErrNotWired) {
		t.Fatalf("SubmitOrder returned err=%v, want ErrNotWired", err)
	}
}

func TestCancelOrder_NotWired(t *testing.T) {
	err := host.CancelOrder("order-1")
	if !errors.Is(err, host.ErrNotWired) {
		t.Fatalf("CancelOrder returned err=%v, want ErrNotWired", err)
	}
}

func TestLevelConstants(t *testing.T) {
	if host.LevelDebug == "" || host.LevelInfo == "" ||
		host.LevelWarn == "" || host.LevelError == "" {
		t.Fatal("all four Level constants must be non-empty")
	}
}
