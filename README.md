# Wazero Buffer Pool

A [wazero](https://pkg.go.dev/github.com/tetratelabs/wazero) host module, ABI and guest SDK providing sets of multi value buffer pools.

## Host Module

[![Go Reference](https://godoc.org/github.com/pantopic/wazero-buffer-pool/host?status.svg)](https://godoc.org/github.com/pantopic/wazero-buffer-pool/host)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-buffer-pool/host)](https://goreportcard.com/report/github.com/pantopic/wazero-buffer-pool/host)
[![Go Coverage](https://github.com/pantopic/wazero-buffer-pool/wiki/host/coverage.svg)](https://raw.githack.com/wiki/pantopic/wazero-buffer-pool/host/coverage.html)

First register the host module with the runtime

```go
import (
    "github.com/tetratelabs/wazero"
    "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

    "github.com/pantopic/wazero-buffer-pool/host"
)

func main() {
    ctx := context.Background()
    r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
    wasi_snapshot_preview1.MustInstantiate(ctx, r)

    module := wazero_buffer_pool.New()
    module.Register(ctx, r)

    // ...
}
```

## Guest SDK (Go)

[![Go Reference](https://godoc.org/github.com/pantopic/wazero-buffer-pool/sdk-go?status.svg)](https://godoc.org/github.com/pantopic/wazero-buffer-pool/sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/wazero-buffer-pool/sdk-go)](https://goreportcard.com/report/github.com/pantopic/wazero-buffer-pool/sdk-go)

Then you can import the guest SDK into your WASI module to send messages from one WASI module to another.

```go
package main

import (
    "github.com/pantopic/wazero-buffer-pool/sdk-go"
)

const (
	BUFFER_POOL_TEST = iota
)

var bpmvs *buffer_pool.MultiValueSet

func main() {
    bpmvs = buffer_pool.NewMutliValueSet(BUFFER_POOL_TEST)
}

//export test
func test() {
    buf := bpmvs.Find(1)
    buf.Append([]byte(`a`))
    buf.Append([]byte(`b`))
    for val := range buf.Iter() {
        println(string(val)) // a, b
    }
    buf.Reset()
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
