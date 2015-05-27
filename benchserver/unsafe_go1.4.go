// +build go1.4,!go1.5

package benchserver

import (
	"runtime"
	"sync"
	"testing"
	"unsafe"
)

// the function doesn't actually return anything, we just need it to
// return something to make dependency tracking working.
func init14() bool

var inited = init14()

var tests *[]testing.InternalTest
var benchmarks *[]InternalBenchmark // tricky, switched type

var benchmarkLock sync.Mutex

// _runN is copied verbatim from the Go 1.4 testing package ((*testing.B).runN).

// _runN runs a single benchmark for the specified number of iterations.
func _runN(b *B, n int) {
	// always safe, we've checked the types are identical
	b2 := (*testing.B)(unsafe.Pointer(b))
	benchmarkLock.Lock()
	defer benchmarkLock.Unlock()
	// Try to get a comparable environment for each run
	// by clearing garbage from previous runs.
	runtime.GC()
	b.N = n
	b.parallelism = 1
	b2.ResetTimer()
	b2.StartTimer()
	b.benchmark.F(b)
	b2.StopTimer()
	b.previousN = n
	b.previousDuration = b.duration
}

func runN(b *B, n int) *testing.BenchmarkResult {
	_runN(b, n)
	// see (*testing.B).run
	return &testing.BenchmarkResult{b.N, b.duration, b.bytes, b.netAllocs, b.netBytes}
}
