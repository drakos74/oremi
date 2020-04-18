package model

import (
	"log"
	"math"
)

// Vector defines a point in n dimensional space
type Vector struct {
	Label  string
	Coords []float64
}

// NewVector creates a new point at the specified coordinates
func NewVector(label string, x ...float64) Vector {
	return Vector{Label: label, Coords: x}
}

// Norm returns the norm of the point
// e.g. the distance to the start of the coordinate system
func (p Vector) Norm() float64 {
	var x2 float64
	for _, x := range p.Coords {
		x2 += math.Pow(x, 2)
	}
	return math.Sqrt(x2)
}

// Dim returns the dimensions of the point
func (p Vector) Dim() int {
	return len(p.Coords)
}

// Distance calculates the distance of the point to another vector
func (p Vector) Distance(element Vector) float64 {
	dimensionValidator(p, element)
	var d float64
	for i, x := range p.Coords {
		d += math.Pow(x-p.Coords[i], 2)
	}
	return math.Sqrt(d)
}

// Iterator iterates over all vectors it is used to read the records of a collection
type Iterator interface {
	Next() (vector Vector, ok, hasNext bool)
	Reset()
}

type Collection interface {
	Iterator
	Size() int
	Edge() (min, max Vector)
}

// Collection represents a collection of vectors e.g. a graph element as such
type CollectionBuilder interface {
	Add(vector Vector)
	Create() Collection
}

var dimensionValidator = func(e1, e2 Vector) {
	if e1.Dim() != e2.Dim() {
		log.Fatalf("vectors don't have the same dimension %v vs %v", e1, e2)
	}
}
