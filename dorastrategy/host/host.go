// Package host is the host-call surface a generated WASM plugin uses.
// Each function here is backed by a wazero import that the server-side
// runtime wires in Plan 2. Until that wiring lands, the functions
// return ErrNotWired so a plugin compiled against this surface can be
// smoke-tested in isolation.
package host

import (
	"errors"
	"time"
)

// ErrNotWired is returned by every host-call function until the wazero
// runtime in Plan 2 wires them to the server-side implementations. A
// plugin that hits this in production is misconfigured: the validate
// smoke should have caught it.
var ErrNotWired = errors.New("host: not wired (Plan 2)")

// Level is the severity of a log line.
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// OrderIntent is the host-side mirror of dorastrategy.OrderIntent. We
// re-declare it here to avoid an import cycle between the framework
// types and the wazero import surface. The fields mirror Dora's
// CreateOrderRequest (see dora-client-go and the dora-api OpenAPI
// spec); numeric fields are decimal strings to round-trip the
// wire format without precision loss.
type OrderIntent struct {
	Side               string
	Quantity           string
	Type               string
	Price              string
	InverseLeverage    string
	FromGlobalPosition bool
}

// Log emits a structured log line. Key-value pairs are host-side;
// the host decides how to format them.
func Log(level Level, msg string, kv ...string) error {
	return ErrNotWired
}

// GetParamString returns a runtime parameter declared in the
// strategy's manifest.params_schema. Returns ErrParamMissing if the key
// is not declared; ErrParamType if the declared type is not "string".
func GetParamString(key string) (string, error) {
	return "", ErrNotWired
}

// GetParamInt returns an int64 runtime parameter.
func GetParamInt(key string) (int64, error) {
	return 0, ErrNotWired
}

// GetParamFloat returns a float64 runtime parameter.
func GetParamFloat(key string) (float64, error) {
	return 0, ErrNotWired
}

// GetParamBool returns a bool runtime parameter.
func GetParamBool(key string) (bool, error) {
	return false, ErrNotWired
}

// SubmitOrder sends an order intent to the server's order broker. The
// call is synchronous: the plugin blocks until the order is acked
// (orderID returned) or denied (typed error returned).
func SubmitOrder(intent OrderIntent) (string, error) {
	return "", ErrNotWired
}

// CancelOrder cancels a previously submitted order. Returns nil on
// success, typed error on denial.
func CancelOrder(orderID string) error {
	return ErrNotWired
}

// Now returns the host's wall clock. TinyGo doesn't have a stable
// wall clock in WASI preview 1; the host supplies one. Implemented in
// the escape package (re-exported here for plugin-side convenience).
func Now() time.Time {
	return time.Time{}
}
