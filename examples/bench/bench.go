package main

import (
	"flag"
	"fmt"
	"github/drakos74/oremi/bench"
	oremi "github/drakos74/oremi/internal"
	"github/drakos74/oremi/internal/data/model"
	"log"

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

	oremi.DrawScene("benchmarks", layout.Horizontal, 1400, 800,
		map[string][]model.Collection{
			"cpu":    {benchmarks.Extract(bench.Operations, bench.Latency)},
			"memory": {benchmarks.Extract(bench.Heap, bench.Throughput)},
		},
	)
}
