package main

import (
	"encoding/binary"

	"github.com/pantopic/wazero-buffer-pool/sdk-go"
)

const (
	BUFFER_POOL_MULTI_SET_1 = iota
)

var (
	testMulti buffer_pool.MultiValueSet
)

func main() {
	testMulti = buffer_pool.NewMultiValueSet(BUFFER_POOL_MULTI_SET_1)
}

//export testMultiSetAppend
func testMultiSetAppend(id, v uint64) {
	testMulti.Find(id).Append(binary.LittleEndian.AppendUint64([]byte{}, v))
}

//export testMultiSetIter
func testMultiSetIter(id uint64) (total uint64) {
	for item := range testMulti.Find(id).Iter() {
		total += binary.LittleEndian.Uint64(item)
	}
	return
}

//export testMultiSetReset
func testMultiSetReset(id uint64) {
	testMulti.Find(id).Reset()
}
