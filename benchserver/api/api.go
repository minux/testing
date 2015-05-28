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
