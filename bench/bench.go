package bench

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github/drakos74/oremi/internal/data/model"
)

const (
	operations = "ops"
	latency    = "ns/op"
	throughput = "B/op"
	heap       = "allocs/op"
)

// NewBenchmarkCollection creates a benchmark collection of data points from a benchmark output file
// x - axis latency
// y - axis operations
func NewBenchmarkCollection(f string) (model.Collection, error) {

	series := model.NewSeries(2)

	file, err := os.Open(f)
	if err != nil {
		return series, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		b, err := tryParseBenchmark(line)
		if err == nil {
			series.Add(model.NewVector(fmt.Sprintf("%v", b.labels), b.Latency(), b.Operations()))
		} else {
			println(fmt.Sprintf("ignoring line %s : %v", line, err))
		}
	}

	if err := scanner.Err(); err != nil {
		return series, err
	}

	return series, nil

}

type benchmark struct {
	labels []string
	// numeric labels
	numLabels map[string]float64
}

func (b benchmark) Latency() float64 {
	return b.numLabels[latency]
}

func (b benchmark) Operations() float64 {
	return b.numLabels[operations]
}

func (b benchmark) Throughput() float64 {
	return b.numLabels[throughput]
}

func (b benchmark) Heap() float64 {
	return b.numLabels[heap]
}

func tryParseBenchmark(line string) (*benchmark, error) {
	parts := strings.Fields(line)

	if isBenchmarkOutput(parts) {
		return parseBenchmark(parts), nil
	}

	return nil, fmt.Errorf("could not find bench within %v", parts)

}

func parseBenchmark(parts []string) *benchmark {
	ops, _ := strconv.Atoi(parts[1])
	lat, _ := strconv.ParseFloat(parts[2], 64)

	labels, numLabels := parseLabels(parts[0])

	numLabels[latency] = lat
	numLabels[operations] = float64(ops)

	return &benchmark{
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
