//go:build appengine || windows
// +build appengine windows

// base source: https://github.com/VictoriaMetrics/fastcache/blob/master/malloc_heap.go
package malloc

func getChunk(chunkSize int) []byte {
	return make([]byte, chunkSize)
}

func putChunk(chunk []byte) {
	// No-op.
}
