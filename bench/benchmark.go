package bench

import (
	"bufio"
	"fmt"
	datamodel "github/drakos74/oremi/internal/data/model"
	"github/drakos74/oremi/internal/gui"
	"github/drakos74/oremi/internal/gui/canvas/entity"
	uimodel "github/drakos74/oremi/internal/gui/model"
	"log"
	"os"
	"strconv"
	"strings"

	"gioui.org/f32"
)

const (
	operations = "ops"
	latency    = "ns/op"
	throughput = "B/op"
	heap       = "allocs/op"
)

type benchmarks []benchmark

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

// ParseAndPlot will parse a simple benchmark log file and plot the results
// x - axis latency
// y - axis operations
func ParseAndPlot(file string, width, height float32) {
	b, err := parseBenchmarks(file)

	if err != nil {
		log.Fatalf("could not parse benchamrks from file '%s': %v", file, err)
	}

	collection := b.toCollection()

	var scene gui.Scene
	scene.WithDimensions(width, height)

	g := f32.Rectangle{
		Min: f32.Point{X: 50, Y: 50},
		Max: f32.Point{X: width * 2 * 95 / 100, Y: height * 2 * 95 / 100},
	}
	graph := entity.NewGraph(&g)
	graph.AddCollection(uimodel.NewSeries(collection))

	scene.Add(graph)
	scene.Run()
}

func (b benchmarks) toCollection() datamodel.Collection {
	series := datamodel.NewSeries(2)
	for _, benchmark := range b {
		series.Add(datamodel.NewVector(fmt.Sprintf("%v", benchmark.labels), benchmark.Latency(), benchmark.Operations()))
	}
	return series
}

func parseBenchmarks(f string) (benchmarks, error) {

	var benchmarks []benchmark

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

func parseBenchmark(line string) (*benchmark, error) {
	parts := strings.Fields(line)

	if isBenchmarkOutput(parts) {
		return parseAsBenchmark(parts), nil
	}

	return nil, fmt.Errorf("could not find bench within %v", parts)

}

func parseAsBenchmark(parts []string) *benchmark {
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
