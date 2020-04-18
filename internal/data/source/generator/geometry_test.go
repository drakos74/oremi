package generator

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLine(t *testing.T) {

	line := NewLine(10, 2, 1, 1)

	assert.Equal(t, 10, line.Size())

	min, max := line.Edge()

	assert.NotEqual(t, min, max)

	assert.Equal(t, float64(0), math.Round(min.Coords[0]))
	assert.Equal(t, float64(1), math.Round(min.Coords[1]))

	assert.Equal(t, float64(1+(10-1)*2), math.Round(max.Coords[1]))

	println(fmt.Sprintf("line = %v", line))

}

func TestNewPolynomial(t *testing.T) {

	poly := NewPolynomial(100, 0, 0.1, 0, 1)

	assert.Equal(t, 100, poly.Size())

	min, max := poly.Edge()

	assert.NotEqual(t, min, max)

	assert.Equal(t, float64(0), math.Round(min.Coords[0]))
	assert.Equal(t, float64(0), math.Round(min.Coords[1]))

	assert.Equal(t, math.Round(math.Pow(10-0.1, 2)), math.Round(max.Coords[1]))

	println(fmt.Sprintf("poly = %v", poly))

}

func TestNewExponential(t *testing.T) {

	exp := NewExponential(100, 1, 0.1)

	assert.Equal(t, 100, exp.Size())

	min, max := exp.Edge()

	assert.NotEqual(t, min, max)

	assert.Equal(t, float64(0), math.Round(min.Coords[0]))
	assert.Equal(t, float64(1), math.Round(min.Coords[1]))

	assert.Equal(t, math.Round(math.Exp(10-0.1)), math.Round(max.Coords[1]))

	println(fmt.Sprintf("exp = %v", exp))

}
