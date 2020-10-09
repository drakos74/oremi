package bench

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkFilter(b *testing.B) {

	for i := 0; i < 25; i++ {

		n := i * 2
		s := i

		filterLabels, benchmarkLabels := generateNumLabels(i*2, i)
		benchmark := Benchmark{
			numLabels: benchmarkLabels,
		}

		includeFilter := Include(filterLabels)
		excludeFilter := Exclude(filterLabels)

		match := true

		b.Run(Format(map[string]int{
			"keys":   n,
			"length": s,
		}, "include"), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				includeFilter.Apply(benchmark, &match)
			}
		})

		b.Run(Format(map[string]int{
			"keys":   n,
			"length": s,
		}, "exclude"), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				excludeFilter.Apply(benchmark, &match)
			}
		})
	}

}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateNumLabels(n, s int) (filter, bench map[string]float64) {
	filter = map[string]float64{}
	bench = map[string]float64{}
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		b := make([]byte, s)
		for c := range b {
			b[c] = charset[seed.Intn(len(charset))]
		}
		key := string(b)
		value := seed.Float64()
		otherValue := seed.Float64()
		filter[key] = value
		bench[key] = otherValue

	}
	return
}
