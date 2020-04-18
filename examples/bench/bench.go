package main

import (
	"flag"
	"fmt"
	"github/drakos74/oremi/bench"
	"log"
)

func main() {

	file := flag.String("file", "examples/bench/testdata/benchmark_output.txt", "bench output file")

	flag.Parse()

	println(fmt.Sprintf("parsing bench file = %v", *file))

	collection, err := bench.NewBenchmarkCollection(*file)

	if err != nil {
		log.Fatalf("could not parse benchamrks from file '%s': %v", file, err)
	}

	bench.DrawCollection(collection, 1200, 800)

}
