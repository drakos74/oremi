package bench

import (
	"bufio"
	"fmt"
	"github/drakos74/oremi/internal/data/model"
	"os"
	"strconv"
	"strings"
)

const (
	operations = "ops"
	latency    = "ns/op"
	throughput = "B/op"
	heap       = "allocs/op"
)

type Benchmarks []Benchmark

type Benchmark struct {
	labels []string
	// numeric labels
	numLabels map[string]float64
}

func (b Benchmark) Latency() float64 {
	return b.numLabels[latency]
}

func (b Benchmark) Operations() float64 {
	return b.numLabels[operations]
}

func (b Benchmark) Throughput() float64 {
	return b.numLabels[throughput]
}

func (b Benchmark) Heap() float64 {
	return b.numLabels[heap]
}

func (b Benchmarks) ToCollection() model.Collection {
	series := model.NewSeries(2)
	for _, benchmark := range b {
		series.Add(model.NewVector(fmt.Sprintf("%v", benchmark.labels), benchmark.Latency(), benchmark.Operations()))
	}
	return series
}

func ParseBenchmarks(f string) (Benchmarks, error) {

	var benchmarks []Benchmark

	file, err := os.Open(f)
	if err != nil {
		return benchmarks, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		b, err := parseBenchmark(line)
		if err == nil {
			benchmarks = append(benchmarks, *b)
		} else {
			println(fmt.Sprintf("ignoring line %s : %v", line, err))
		}
	}

	if err := scanner.Err(); err != nil {
		return benchmarks, err
	}

	return benchmarks, nil

}

func parseBenchmark(line string) (*Benchmark, error) {
	parts := strings.Fields(line)

	if isBenchmarkOutput(parts) {
		return parseAsBenchmark(parts), nil
	}

	return nil, fmt.Errorf("could not find benchmark within %v", parts)

}

func parseAsBenchmark(parts []string) *Benchmark {
	ops, _ := strconv.Atoi(parts[1])
	lat, _ := strconv.ParseFloat(parts[2], 64)

	labels, numLabels := parseLabels(parts[0])

	numLabels[latency] = lat
	numLabels[operations] = float64(ops)

	return &Benchmark{
		labels:    labels,
		numLabels: numLabels,
	}
}

func parseLabels(str string) (labels []string, numLabels map[string]float64) {

	labels = make([]string, 0)
	numLabels = make(map[string]float64)

	parts := strings.Split(str, "|")

	for _, p := range parts {
		label := strings.Split(p, ":")
		if len(label) > 1 {
			num, _ := strconv.ParseFloat(label[1], 64)
			numLabels[label[0]] = num
		} else {
			labels = append(labels, label[0])
		}

	}

	return

}

func isBenchmarkOutput(parts []string) bool {
	return hasAtIndex(3, "ns/op", parts)
}

func hasAtIndex(index int, part string, parts []string) bool {
	if len(parts) > index {
		return parts[index] == part
	}
	return false
}
