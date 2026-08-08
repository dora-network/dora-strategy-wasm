# WASM framework sub-module build helpers.
#
# The default `go build ./...` covers the framework packages. The
# `tinygo-build` target is the smoke test for the WASM pipeline:
# it compiles the example strategy to a .wasm module. CI skips
# this target if `tinygo` is not on PATH.

GO ?= go
TINYGO ?= tinygo

.PHONY: build
build:
	$(GO) build ./...

.PHONY: test
test:
	$(GO) test ./...

.PHONY: tinygo-build
tinygo-build:
	@if ! command -v $(TINYGO) >/dev/null 2>&1; then \
		echo "tinygo not on PATH; skipping (set TINYGO=path/to/binary to override)"; \
		exit 0; \
	fi
	$(TINYGO) build -target=wasi -buildmode=c-shared -o example/strategy.wasm example/strategy.go
	@head -c 4 example/strategy.wasm | od -An -c | grep -q '\\0   a   s   m' && \
		echo "tinygo-build: ok (wasm magic bytes present)" || \
		(echo "tinygo-build: FAILED (no wasm magic bytes)"; exit 1)
