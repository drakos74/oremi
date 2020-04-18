package pkg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBenchmark(t *testing.T) {

	b, err := ParseBenchmark("BenchmarkMemory/*mem.Cache:put/num:0,size-key:2,size-value:10-16            	529767379	         2.21 ns/op")
	assert.NoError(t, err)

	assert.Equal(t, 2.21, b.latency)

	println(fmt.Sprintf("%v", b))
}

func TestParseBenchmarks(t *testing.T) {

	benchmarks, err := ParseBenchmarks("benchmark_output.txt")
	assert.NoError(t, err)

	for _, b := range benchmarks {
		println(fmt.Sprintf("%v", b))
	}
}
