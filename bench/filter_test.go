package bench

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type scenario struct {
	filterLabels []string
	filter       map[string]float64
	labels       []string
	numLabels    map[string]float64
	prevMatch    bool
	result       bool
}

var includeScenarios = []scenario{
	{
		// no match in any property
		filter:    map[string]float64{"test": 1, "demo": 2},
		numLabels: map[string]float64{"test": 10, "demo": 20},
		prevMatch: true,
		result:    false,
	},
	{
		// match in one property
		filter:    map[string]float64{"test": 1, "demo": 20},
		numLabels: map[string]float64{"test": 10, "demo": 20},
		prevMatch: true,
		result:    false,
	},
	{
		// match in both properties
		filter:    map[string]float64{"test": 10, "demo": 20},
		numLabels: map[string]float64{"test": 10, "demo": 20},
		prevMatch: true,
		result:    true,
	},
	{
		// match in both properties with previous result being false
		filter:    map[string]float64{"test": 10, "demo": 20},
		numLabels: map[string]float64{"test": 10, "demo": 20},
		prevMatch: false,
		result:    false,
	},
}

func TestInclude(t *testing.T) {

	for _, testCase := range includeScenarios {
		filter := Include(testCase.filter)

		benchmark := Benchmark{
			numLabels: testCase.numLabels,
		}

		match := testCase.prevMatch
		filter.Apply(benchmark, &match)
		assert.Equal(t, testCase.result, match)
	}

}

var excludeScenarios = []scenario{
	{
		// no match in any property
		filter:    map[string]float64{"test": 1, "demo": 2},
		numLabels: map[string]float64{"test": 10, "demo": 20},
		prevMatch: true,
		result:    true,
	},
	{
		// match in one property
		filter:    map[string]float64{"test": 1, "demo": 20},
		numLabels: map[string]float64{"test": 10, "demo": 20},
		prevMatch: true,
		result:    false,
	},
	{
		// match in both properties
		filter:    map[string]float64{"test": 10, "demo": 20},
		numLabels: map[string]float64{"test": 10, "demo": 20},
		prevMatch: true,
		result:    false,
	},
	{
		// match in both properties with previous result being false
		filter:    map[string]float64{"test": 10, "demo": 20},
		numLabels: map[string]float64{"test": 10, "demo": 20},
		prevMatch: false,
		result:    false,
	},
}

func TestExclude(t *testing.T) {

	for _, testCase := range excludeScenarios {
		filter := Exclude(testCase.filter)

		benchmark := Benchmark{
			numLabels: testCase.numLabels,
		}

		match := testCase.prevMatch
		filter.Apply(benchmark, &match)
		assert.Equal(t, testCase.result, match)
	}

}

var labelScenarios = []scenario{
	{
		// no match in any property
		filterLabels: []string{"other", "another"},
		labels:       []string{"test", "demo"},
		prevMatch:    true,
		result:       false,
	},
	{
		// match in one property
		filterLabels: []string{"other", "demo"},
		labels:       []string{"test", "demo"},
		prevMatch:    true,
		result:       false,
	},
	{
		// match in both properties
		filterLabels: []string{"test", "demo"},
		labels:       []string{"test", "demo"},
		prevMatch:    true,
		result:       false,
	},
	{
		// match in both properties with previous result being false
		filterLabels: []string{"test", "demo"},
		labels:       []string{"test", "demo"},
		prevMatch:    false,
		result:       false,
	},
}

func TestLabel(t *testing.T) {

	for _, testCase := range labelScenarios {
		filter := Label(testCase.filterLabels...)

		benchmark := Benchmark{
			numLabels: testCase.numLabels,
		}

		match := testCase.prevMatch
		filter.Apply(benchmark, &match)
		assert.Equal(t, testCase.result, match)
	}

}
