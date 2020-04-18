package main

import (
	"flag"
	"fmt"
	"github/drakos74/oremi/bench"
)

func main() {

	file := flag.String("file", "examples/bench/testdata/benchmark_output.txt", "bench output file")

	flag.Parse()

	println(fmt.Sprintf("parsing bench file = %v", *file))

	bench.ParseAndPlot(*file, 1200, 800)
}
