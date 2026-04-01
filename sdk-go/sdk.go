package buffer_pool

import (
	"encoding/binary"
	"iter"
)

type MultiValueSet struct {
	id   uint64
	size uint64
}

func NewMultiValueSet(id uint64, opts ...Option) MultiValueSet {
	p := MultiValueSet{id: id}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

func (s *MultiValueSet) Find(id uint64) MultiValue {
	return MultiValue{
		id:    id,
		setID: s.id,
		size:  s.size,
	}
}

type MultiValue struct {
	id    uint64
	setID uint64
	size  uint64
}

func (c MultiValue) Append(b []byte) bool {
	id = c.id
	setID = c.setID
	copy(buf[:len(b)], b)
	bufLen = uint32(len(b))
	_multi_append()
	return errCode == 0
}

func (c MultiValue) Iter() iter.Seq[[]byte] {
	id = c.id
	setID = c.setID
	_multi_load()
	return func(yield func([]byte) bool) {
		for i := 0; i < int(bufLen); {
			size, n := binary.Uvarint(buf[i:])
			i += n
			if !yield(buf[i : i+int(size)]) {
				return
			}
			i += int(size)
		}
	}
}

func (c MultiValue) Reset() {
	id = c.id
	setID = c.setID
	_multi_reset()
}
