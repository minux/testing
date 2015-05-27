// +build go1.5

package benchserver

import (
	"testing"
	_ "unsafe" // so that we can use //go:linkname
)

// unsafe hackery for Go 1.5

//go:linkname _tests main.tests
var _tests []testing.InternalTest

//go:linkname _benchmarks main.benchmarks
var _benchmarks []InternalBenchmark // tricky, switched type

var (
	tests      = &_tests
	benchmarks = &_benchmarks
	inited     bool
)

// Here comes the interesting part, although we define runN to
// take *B and int as parameter, it's actually refering to
// the method runN on *testing.B. Essentially we used our B
// to replace testing.B. This is definitely unsafe, so we have
// a test to verify that our definition of B is identical to
// the testing package's.
// We are also assuming that method on type T has the same calling
// convention as a func with an added parameter of type T. This
// is true for both gc and gccgo.

//go:linkname runN testing.(*B).runN
func runN(*B, int)
