package generator

import (
	"fmt"
	"github/drakos74/oremi/internal/data/model"
	"log"
	"math"
)

// Y is a general abstraction of a mathematical function
// Y = f(X)
type Y func(x ...float64) float64

func linearFunction(a, b float64) Y {
	return polynomialFunction(b, a)
}

// y = x0 * x^0 + x1 * x^1 + x2 * x^2 + ...
func polynomialFunction(a ...float64) Y {
	return func(x ...float64) float64 {
		var y float64
		for i, aa := range a {
			y += aa * math.Pow(x[0], float64(i))
		}
		return y
	}
}

// y = a * e^x
func exponentialFunction(a float64) Y {
	return func(x ...float64) float64 {
		return a * math.Exp(x[0])
	}
}

// F is a collection of functions that transform the respective coordinates
// for the sake of the abstraction , we do incude all coordinates
// e.g. [ y1,y2,y3...] = [ f1(x1,x2,x3...) , f2(x1,x2,x3...) , f3(x1,x2,x3...) ]
type F struct {
	f []Y
}

// Call executes the F function on a vector
func (g *F) Call(p ...float64) []float64 {
	coords := make([]float64, len(p))
	for i := 0; i < len(p); i++ {
		coords[i] = g.f[i](p...)
	}
	return coords
}

// GN represents a generator to produce a model series with the specified number of objects
type GN interface {
	Num(num int) model.Collection
}

// GL represents a generator to produce a model series with the Limit
type GL interface {
	Lim(limit float64) model.Collection
}

type Euclidean struct {
	model.Iterator
	F
}

// Generate generates a new number of points for the graph series
func (g Euclidean) Num(num int) model.Collection {

	l := model.NewSeries("x", "f(x)")

	for i := 0; i < num; i++ {

		if x, ok, n := g.Next(); ok {
			p := g.Call(x.Coords[0], 0)
			l.Add(model.NewVector([]string{"f", fmt.Sprintf("%d", i)}, p...))
			if !n {
				log.Printf("could not generate more elements, source iterator ended at %d of %d", i, num)
				break
			}
		}

	}
	return l
}

// Generate generates a new number of points for the graph series
func (g Euclidean) Lim(limit float64) model.Collection {

	l := model.NewSeries("x", "f(x)")

	i := 0

	for {
		if x, ok, n := g.Next(); ok {
			p := g.Call(x.Coords[0], 0)
			l.Add(model.NewVector([]string{"f", fmt.Sprintf("%d", i)}, p...))
			if !n {
				log.Printf("cannot not generate more elements, source iterator ended at %d", i)
				break
			}
			if p[1] > limit {
				log.Printf("cannot not generate more elements, source iterator ended at %f", p[1])
				break
			}
		}
		i++
	}
	return l
}
