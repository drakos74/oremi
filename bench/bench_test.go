package bench

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBenchmark(t *testing.T) {

	b, err := tryParseBenchmark("BenchmarkMemory/*mem.Cache|put|num:0|size-key:2|size-value:10|-16            	529767379	         2.21 ns/op")
	assert.NoError(t, err)

	latency, ok := b.Latency()
	assert.True(t, ok)
	assert.Equal(t, 2.21, latency)

	assert.Equal(t, []string{"BenchmarkMemory/*mem.Cache", "put", "-16"}, b.labels)
	assert.Equal(t, map[string]float64{Operations: float64(529767379), Latency: float64(2.21), "num": 0, "size-key": 2, "size-value": 10}, b.numLabels)
}

func TestParseBenchmarkWithAllocs(t *testing.T) {
	b, err := tryParseBenchmark("BenchmarkSB/*file.SB|Get|num:1000|size-key:4|size-value:100|-16                      608           1945741 ns/op          208000 B/op       3000 allocs/op")

	assert.NoError(t, err)

	latency, ok := b.Latency()
	assert.True(t, ok)
	assert.Equal(t, float64(1945741), latency)
	throughput, ok := b.Throughput()
	assert.True(t, ok)
	assert.Equal(t, float64(208000), throughput)
	heap, ok := b.Heap()
	assert.True(t, ok)
	assert.Equal(t, float64(3000), heap)

	assert.Equal(t, []string{"BenchmarkSB/*file.SB", "Get", "-16"}, b.labels)
	assert.Equal(t, map[string]float64{Throughput: float64(208000), Heap: float64(3000), Operations: float64(608), Latency: float64(1945741), "num": 1000, "size-key": 4, "size-value": 100}, b.numLabels)
}

func TestParseBenchmarks(t *testing.T) {

	benchmarks, err := New("bench_output/filter_example_benchmark_test.txt")
	assert.NoError(t, err)

	assert.Equal(t, 50, len(benchmarks))

	for _, b := range benchmarks {
		println(fmt.Sprintf("%v", b))
	}
}
