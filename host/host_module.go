package wazero_buffer_pool

import (
	"context"
	"encoding/binary"
	"log"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Name is the name of this host module.
const Name = "pantopic/wazero-buffer-pool"

var (
	ctxKeyMeta = Name + `/meta`
	ctxKeyPool = Name + `/pool`
)

type meta struct {
	ptrBuf     uint32
	ptrBufCap  uint32
	ptrBufLen  uint32
	ptrErrCode uint32
	ptrID      uint32
	ptrSetID   uint32
}

type hostModule struct {
	sync.RWMutex

	module api.Module
}

type Option func(*hostModule)

func New(opts ...Option) *hostModule {
	p := &hostModule{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (h *hostModule) Name() string {
	return Name
}

func (h *hostModule) ContextCopy(dst, src context.Context) context.Context {
	dst = context.WithValue(dst, ctxKeyMeta, get[*meta](src, ctxKeyMeta))
	if v := src.Value(ctxKeyPool); v != nil {
		dst = context.WithValue(dst, ctxKeyPool, v.(map[uint64]map[uint64][]byte))
	} else {
		dst = context.WithValue(dst, ctxKeyPool, make(map[uint64]map[uint64][]byte))
	}
	return dst
}

func (h *hostModule) Stop() {}

var scratch = make([]byte, 8)

// Register instantiates the host module, making it available to all module instances in this runtime
func (h *hostModule) Register(ctx context.Context, r wazero.Runtime) (err error) {
	builder := r.NewHostModuleBuilder(Name)
	register := func(name string, fn func(ctx context.Context, m api.Module, stack []uint64)) {
		builder = builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(fn), nil, nil).Export(name)
	}
	for name, fn := range map[string]any{
		"__buffer_pool_multi_append": func(m map[uint64][]byte, id uint64, v []byte) bool {
			if _, ok := m[id]; !ok {
				m[id] = []byte{}
			}
			n := binary.PutUvarint(scratch, uint64(len(v)))
			if len(m[id])+len(v)+n > cap(v) {
				return false
			}
			m[id] = append(m[id], scratch[:n]...)
			m[id] = append(m[id], v...)
			return true
		},
		"__buffer_pool_multi_load": func(m map[uint64][]byte, id uint64) []byte {
			return m[id]
		},
		"__buffer_pool_multi_reset": func(m map[uint64][]byte, id uint64) {
			// TODO - pool buffers instead of deleting
			delete(m, id)
		},
	} {
		switch fn := fn.(type) {
		case func(m map[uint64][]byte, id uint64, v []byte) bool:
			register(name, func(ctx context.Context, mod api.Module, stack []uint64) {
				meta := get[*meta](ctx, ctxKeyMeta)
				errCode := uint32(0)
				if !fn(h.getMap(ctx, mod, meta), getID(mod, meta), getBuf(mod, meta)) {
					errCode = 1
				}
				writeUint32(mod, meta.ptrErrCode, errCode)
			})
		case func(m map[uint64][]byte, id uint64) []byte:
			register(name, func(ctx context.Context, mod api.Module, stack []uint64) {
				meta := get[*meta](ctx, ctxKeyMeta)
				b := fn(h.getMap(ctx, mod, meta), getID(mod, meta))
				copy(buf(mod, meta)[:len(b)], b)
				writeUint32(mod, meta.ptrBufLen, uint32(len(b)))
			})
		case func(m map[uint64][]byte, id uint64):
			register(name, func(ctx context.Context, mod api.Module, stack []uint64) {
				meta := get[*meta](ctx, ctxKeyMeta)
				fn(h.getMap(ctx, mod, meta), getID(mod, meta))
			})
		default:
			log.Panicf("Method signature implementation missing: %#v", fn)
		}
	}
	h.module, err = builder.Instantiate(ctx)
	return
}

// InitContext retrieves the meta page from the wasm module
func (h *hostModule) InitContext(ctx context.Context, m api.Module) (context.Context, error) {
	stack, err := m.ExportedFunction(`__buffer_pool`).Call(ctx)
	if err != nil {
		return ctx, err
	}
	meta := &meta{}
	ptr := uint32(stack[0])
	for i, v := range []*uint32{
		&meta.ptrID,
		&meta.ptrSetID,
		&meta.ptrBufCap,
		&meta.ptrBufLen,
		&meta.ptrBuf,
		&meta.ptrErrCode,
	} {
		*v = readUint32(m, ptr+uint32(4*i))
	}
	return context.WithValue(ctx, ctxKeyMeta, meta), nil
}

func (h *hostModule) getMap(ctx context.Context, mod api.Module, meta *meta) map[uint64][]byte {
	id := readUint64(mod, meta.ptrID)
	m := get[map[uint64]map[uint64][]byte](ctx, ctxKeyPool)
	h.RLock()
	_, ok := m[id]
	h.RUnlock()
	if !ok {
		h.Lock()
		if _, ok := m[id]; !ok {
			m[id] = map[uint64][]byte{}
		}
		h.Unlock()
	}
	return m[id]
}

func getID(mod api.Module, meta *meta) uint64 {
	return readUint64(mod, meta.ptrID)
}

func getBuf(mod api.Module, meta *meta) []byte {
	return read(mod, meta.ptrBuf, meta.ptrBufLen, meta.ptrBufCap)
}

func getBufCopy(mod api.Module, meta *meta) []byte {
	return append([]byte(nil), getBuf(mod, meta)...)
}

func buf(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrBuf, 0, meta.ptrBufCap)
}

func get[T any](ctx context.Context, key string) T {
	v := ctx.Value(key)
	if v == nil {
		log.Panicf("Context item missing %s", key)
	}
	return v.(T)
}

func id(m api.Module, meta *meta) uint32 {
	return readUint32(m, meta.ptrID)
}

func readUint32(m api.Module, ptr uint32) (val uint32) {
	val, ok := m.Memory().ReadUint32Le(ptr)
	if !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
	return
}

func read(m api.Module, ptrData, ptrLen, ptrCap uint32) (buf []byte) {
	buf, ok := m.Memory().Read(ptrData, readUint32(m, ptrCap))
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", ptrData, ptrLen)
	}
	return buf[:readUint32(m, ptrLen)]
}

func readUint64(m api.Module, ptr uint32) (val uint64) {
	val, ok := m.Memory().ReadUint64Le(ptr)
	if !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
	return
}

func writeUint64(m api.Module, ptr uint32, val uint64) {
	if ok := m.Memory().WriteUint64Le(ptr, val); !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
}

func writeUint32(m api.Module, ptr uint32, val uint32) {
	if ok := m.Memory().WriteUint32Le(ptr, val); !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
}
