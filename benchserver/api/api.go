// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"testing"
)

type Arg struct {
	Name string `json:"name"`
	N    int    `json:"n"`
}

type Reply struct {
	Result *testing.BenchmarkResult `json:"result,omitempty"`
	Names  []string                 `json:"names,omitempty"`
}
