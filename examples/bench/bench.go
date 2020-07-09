package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/drakos74/oremi"

	"github.com/drakos74/oremi/bench"

	"gioui.org/layout"
)

func main() {

	file := flag.String("file", "examples/bench/data/benchmark_output_ext.txt", "bench output file")

	flag.Parse()

	println(fmt.Sprintf("parsing benchmark file = %v", *file))

	benchmarks, err := bench.New(*file)

	if err != nil {
		log.Fatalf("could not parse benchamrks from file '%s': %v", *file, err)
	}

	oremi.Draw("benchmarks", layout.Horizontal, 1400, 800, gatherBenchmarks(benchmarks))

}

func gatherBenchmarks(benchmarks bench.Benchmarks) map[string]map[string]oremi.Collection {

	graphs := make(map[string]bench.Benchmarks)
	colors := bench.Palette()

	for _, b := range benchmarks {
		println(fmt.Sprintf("b = %v", b))
		label := b.Labels()[0]
		i := strings.Index(label, "/")
		l := label[0:i]
		if _, ok := graphs[l]; !ok {
			graphs[l] = make([]bench.Benchmark, 0)
		}
		graphs[l] = append(graphs[l], b)
	}

	collections := make(map[string]map[string]oremi.Collection)
	collections["latency"] = make(map[string]oremi.Collection)
	//collections["memory"] = make(map[string]oremi.Collection)
	for label, graph := range graphs {
		collections["latency"][label] = graph.Extract(bench.Operations, bench.Latency, map[string]float64{bench.Num: 1000}).
			Color(colors.Get(label))
		//collections["memory"][label] = graph.Extract(bench.Heap, bench.Throughput).
		//	Color(colors.Get(label))
	}
	return collections
}
