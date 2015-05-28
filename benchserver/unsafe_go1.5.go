// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// Here comes the interesting part, although we define _runN to
// take *B and int as parameter, it's actually refering to
// the method runN on *testing.B. Essentially we used our B
// to replace testing.B. This is definitely unsafe, so we have
// a test to verify that our definition of B is identical to
// the testing package's.
// We are also assuming that method on type T has the same calling
// convention as a func with an added parameter of type T. This
// is true for both gc and gccgo.

//go:linkname _runN testing.(*B).runN
func _runN(*B, int)

func runN(b *B, n int) *testing.BenchmarkResult {
	_runN(b, n)
	// see (*testing.B).run
	return &testing.BenchmarkResult{b.N, b.duration, b.bytes, b.netAllocs, b.netBytes}
}
