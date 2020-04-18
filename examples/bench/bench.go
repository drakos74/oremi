package main

import (
	"flag"
	"fmt"
	"github/drakos74/oremi/bench"
	"log"
)

func main() {

	file := flag.String("file", "examples/bench/testdata/benchmark_output_ext.txt", "bench output file")

	flag.Parse()

	println(fmt.Sprintf("parsing bench file = %v", *file))

	benchmarks, err := bench.New(*file)

	if err != nil {
		log.Fatalf("could not parse benchamrks from file '%s': %v", file, err)
	}

	bench.DrawCollections(1600, 800,
		benchmarks.Extract(bench.Operations, bench.Latency),
		benchmarks.Extract(bench.Heap, bench.Throughput),
	)
}
