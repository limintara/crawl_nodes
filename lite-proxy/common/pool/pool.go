package pool

import (
	"bytes"
	"sync"
)

const (
	// io.Copy default buffer size is 32 KiB
	// but the maximum packet size of vmess/shadowsocks is about 16 KiB
	// so define a buffer of 20 KiB to reduce the memory of each TCP relay
	RelayBufferSize = 20 * 1024
)

func Get(size int) []byte {
	return defaultAllocator.Get(size)
}

func Put(buf []byte) error {
	return defaultAllocator.Put(buf)
}

var bufferPool = sync.Pool{New: func() interface{} { return &bytes.Buffer{} }}

// FIXME
func GetBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}
