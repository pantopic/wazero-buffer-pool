package buffer_pool

import (
	"unsafe"
)

var (
	id      uint64
	setID   uint64
	bufCap  uint32 = 1536 << 10 // 1.5 MiB
	bufLen  uint32
	buf     = make([]byte, bufCap)
	errCode uint32
	meta    = make([]uint32, 6)
)

//export __buffer_pool
func __buffer_pool() (res uint32) {
	for i, p := range []unsafe.Pointer{
		unsafe.Pointer(&id),
		unsafe.Pointer(&setID),
		unsafe.Pointer(&bufCap),
		unsafe.Pointer(&bufLen),
		unsafe.Pointer(&buf[0]),
		unsafe.Pointer(&errCode),
	} {
		meta[i] = uint32(uintptr(p))
	}
	return uint32(uintptr(unsafe.Pointer(&meta[0])))
}

//go:wasm-module pantopic/wazero-buffer-pool
//export __buffer_pool_multi_append
func _multi_append()

//go:wasm-module pantopic/wazero-buffer-pool
//export __buffer_pool_multi_load
func _multi_load()

//go:wasm-module pantopic/wazero-buffer-pool
//export __buffer_pool_multi_reset
func _multi_reset()

// Fix for lint rule `unusedfunc`
var _ = __buffer_pool
