package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Benchmark struct {
	name       string
	operations int
	// ns/op
	latency float64
	// B/op
	throughput float64
	// allocs/op
	heap float64
}

func ParseBenchmarks(f string) ([]Benchmark, error) {

	var benchmarks []Benchmark

	file, err := os.Open(f)
	if err != nil {
		return benchmarks, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		b, err := ParseBenchmark(line)
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

func ParseBenchmark(line string) (*Benchmark, error) {
	parts := strings.Fields(line)

	if isBenchmarkOutput(parts) {
		return parseAsBenchmark(parts), nil
	}

	return nil, fmt.Errorf("could not find benchmark within %v", parts)

}

func parseAsBenchmark(parts []string) *Benchmark {
	ops, _ := strconv.Atoi(parts[1])
	lat, _ := strconv.ParseFloat(parts[2], 64)
	return &Benchmark{
		name:       parts[0],
		operations: ops,
		latency:    lat,
	}
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
