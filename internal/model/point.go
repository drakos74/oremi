package model

import (
	"math"
)

type Point struct {
	label  string
	coords []float64
}

func NewPoint(label string, x ...float64) *Point {
	return &Point{label: label, coords: x}
}

func (p Point) Norm() float64 {
	var x2 float64
	for _, x := range p.coords {
		x2 += math.Pow(x, 2)
	}
	return math.Sqrt(x2)
}

func (p Point) Coords() []float64 {
	return p.coords
}

func (p Point) Dim() int {
	return len(p.coords)
}

func (p Point) Distance(element Element) float64 {
	dimensionValidator(p, element)
	var d float64
	for i, x := range p.coords {
		d += math.Pow(x-element.Coords()[i], 2)
	}
	return math.Sqrt(d)
}

func (p Point) Label() string {
	return p.label
}
