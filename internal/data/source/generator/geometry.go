package generator

import (
	"github/drakos74/oremi/internal/data/model"
)

// Line is a collection of points adhering to a graph evolution model in 2D space
type Line model.Series

// NewLine creates a new collection of a number of vectors forming a line
func NewLine(num int, a, b, step float64) model.Collection {
	line := NewLineGenerator(a, b, step)
	collection := line.Num(num)
	return collection
}

// NewLineGenerator creates a new generator for a 2d line
func NewLineGenerator(a, b, step float64) *Euclidean {
	return &Euclidean{
		NewLinearSequence(0, step),
		F{f: []Y{
			// x coordinate :
			func(x ...float64) float64 {
				return x[0]
			},
			// y coordinate
			linearFunction(a, b),
		}}}
}

// Polynomial is a collection of vectors following a polynomial function
type Polynomial model.Series

// NewLine creates a new collection of vectors forming a line
func NewPolynomial(num int, b, step float64, a ...float64) model.Collection {

	polynomial := NewPolynomialGenerator(b, step, a...)
	collection := polynomial.Num(num)

	return collection
}

// NewLineGenerator creates a new generator for a 2d line
func NewPolynomialGenerator(b, step float64, a ...float64) *Euclidean {
	x := []float64{b}
	x = append(x, a...)
	return &Euclidean{
		NewLinearSequence(0, step),
		F{f: []Y{
			// x coordinate :
			func(x ...float64) float64 {
				return x[0]
			},
			// y coordinate
			polynomialFunction(x...),
		}}}
}

// Exponential is a collection of vectors following an exponential function
type Exponential model.Series

// NewLine creates a new collection of vectors forming a line
func NewExponential(num int, a, step float64) model.Collection {

	exponential := NewExponentialGenerator(a, step)
	collection := exponential.Num(num)

	return collection
}

// NewLineGenerator creates a new generator for a 2d line
func NewExponentialGenerator(a, step float64) *Euclidean {
	return &Euclidean{
		NewLinearSequence(0, step),
		F{f: []Y{
			// x coordinate
			func(x ...float64) float64 {
				return x[0]
			},
			// y coordinate
			exponentialFunction(a),
		}}}
}
