// +build go1.4

package benchserver

import (
	"sync"
	"testing"
	"time"
)

// type InternalBenchmark, common and B are copied from the testing package in verbatim

// An internal type but exported because it is cross-package; part of the implementation
// of the "go test" command.
type InternalBenchmark struct {
	Name string
	F    func(b *B) // tricky! we secretly replaced *testing.B with *B
}

// common holds the elements common between T and B and
// captures common methods such as Errorf.
type common struct {
	mu       sync.RWMutex // guards output and failed
	output   []byte       // Output generated by test or benchmark.
	failed   bool         // Test or benchmark has failed.
	skipped  bool         // Test of benchmark has been skipped.
	finished bool

	start    time.Time // Time test or benchmark started
	duration time.Duration
	self     interface{}      // To be sent on signal channel when done.
	signal   chan interface{} // Output for serial tests.
}

// B is a type passed to Benchmark functions to manage benchmark
// timing and to specify the number of iterations to run.
type B struct {
	common
	N                int
	previousN        int           // number of iterations in the previous run
	previousDuration time.Duration // total duration of the previous run
	benchmark        InternalBenchmark
	bytes            int64
	timerOn          bool
	showAllocResult  bool
	result           testing.BenchmarkResult
	parallelism      int // RunParallel creates parallelism*GOMAXPROCS goroutines
	// The initial states of memStats.Mallocs and memStats.TotalAlloc.
	startAllocs uint64
	startBytes  uint64
	// The net total of this test after being run.
	netAllocs uint64
	netBytes  uint64
}
