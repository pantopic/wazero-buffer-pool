# Wazero Small Cache

A [wazero](https://pkg.go.dev/github.com/tetratelabs/wazero) host module, ABI and guest SDK providing a small cache 
suitable for sharing data between concurrent WASI modules.

## Host Module

[![Go Reference](https://godoc.org/github.com/pantopic/wazero-small-cache/host?status.svg)](https://godoc.org/github.com/pantopic/wazero-small-cache/host)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-small-cache/host)](https://goreportcard.com/report/github.com/pantopic/wazero-small-cache/host)
[![Go Coverage](https://github.com/pantopic/wazero-small-cache/wiki/host/coverage.svg)](https://raw.githack.com/wiki/pantopic/wazero-small-cache/host/coverage.html)

First register the host module with the runtime

```go
import (
    "github.com/tetratelabs/wazero"
    "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

    "github.com/pantopic/wazero-small-cache/host"
)

func main() {
    ctx := context.Background()
    r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
    wasi_snapshot_preview1.MustInstantiate(ctx, r)

    module := wazero_atomic.New()
    module.Register(ctx, r)

    // ...
}
```

## Guest SDK (Go)

[![Go Reference](https://godoc.org/github.com/pantopic/wazero-small-cache/sdk-go?status.svg)](https://godoc.org/github.com/pantopic/wazero-small-cache/sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-small-cache/sdk-go)](https://goreportcard.com/report/github.com/pantopic/wazero-small-cache/sdk-go)

Then you can import the guest SDK into your WASI module to send messages from one WASI module to another.

```go
package main

import (
    "github.com/pantopic/wazero-small-cache/sdk-go"
)

const (
	SMALL_CACHE_TEST = iota
)

var c *small_cache.Local

func main() {
    c = small_cache.NewLocal(SMALL_CACHE_TEST)
}

//export test
func test() {
    n.Put([]byte(`a`), []byte(`1`))
    println(n.Get([]byte(`a`))) // 1
}
```

## Roadmap

This project is in alpha. Breaking API changes should be expected until Beta.

- `v0.0.x` - Alpha
  - [ ] Stabilize API
- `v0.x.x` - Beta
  - [ ] Finalize API
  - [ ] Test in production
- `v1.x.x` - General Availability
  - [ ] Proven long term stability in production
