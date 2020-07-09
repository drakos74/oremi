package bench

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBenchmark(t *testing.T) {

	b, err := tryParseBenchmark("BenchmarkMemory/*mem.Cache|put|num:0|size-key:2|size-value:10|-16            	529767379	         2.21 ns/op")
	assert.NoError(t, err)

	assert.Equal(t, 2.21, b.Latency())

	assert.Equal(t, []string{"BenchmarkMemory/*mem.Cache", "put", "-16"}, b.labels)
	assert.Equal(t, map[string]float64{Operations: float64(529767379), Latency: float64(2.21), "num": 0, "size-key": 2, "size-value": 10}, b.numLabels)

	println(fmt.Sprintf("%v", b))
}

func TestParseBenchmarkWithAllocs(t *testing.T) {
	b, err := tryParseBenchmark("BenchmarkSB/*file.SB|Get|num:1000|size-key:4|size-value:100|-16                      608           1945741 ns/op          208000 B/op       3000 allocs/op")

	assert.NoError(t, err)

	assert.Equal(t, float64(1945741), b.Latency())
	assert.Equal(t, float64(208000), b.Throughput())
	assert.Equal(t, float64(3000), b.Heap())

	assert.Equal(t, []string{"BenchmarkSB/*file.SB", "Get", "-16"}, b.labels)
	assert.Equal(t, map[string]float64{Throughput: float64(208000), Heap: float64(3000), Operations: float64(608), Latency: float64(1945741), "num": 1000, "size-key": 4, "size-value": 100}, b.numLabels)

	println(fmt.Sprintf("%v", b))
}

func TestParseBenchmarks(t *testing.T) {

	benchmarks, err := New("test/benchmark_output_ext.txt")
	assert.NoError(t, err)

	assert.Equal(t, 16, len(benchmarks))

	for _, b := range benchmarks {
		println(fmt.Sprintf("%v", b))
	}
}
