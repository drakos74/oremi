package bench

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBenchmark(t *testing.T) {

	b, err := parseBenchmark("BenchmarkMemory/*mem.Cache|put|num:0|size-key:2|size-value:10|-16            	529767379	         2.21 ns/op")
	assert.NoError(t, err)

	assert.Equal(t, 2.21, b.Latency())

	assert.Equal(t, []string{"BenchmarkMemory/*mem.Cache", "put", "-16"}, b.labels)
	assert.Equal(t, map[string]float64{"num": 0, "size-key": 2, "size-value": 10}, b.numLabels)

	println(fmt.Sprintf("%v", b))
}

func TestParseBenchmarks(t *testing.T) {

	benchmarks, err := parseBenchmarks("test/benchmark_output.txt")
	assert.NoError(t, err)

	assert.Equal(t, 16, len(benchmarks))

	for _, b := range benchmarks {
		println(fmt.Sprintf("%v", b))
	}
}
