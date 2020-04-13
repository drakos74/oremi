package model

import (
	"log"
)

type Element interface {
	Norm() float64
	Distance(element Element) float64
	Coords() []float64
	Dim() int
	Label() string
}

var dimensionValidator = func(e1, e2 Element) {
	if e1.Dim() != e2.Dim() {
		log.Fatalf("elements don't have the same dimension %v vs %v", e1, e2)
	}
}

type Iterator interface {
	Next() (element Element, hasNext bool)
}

type Collection interface {
	Iterator
	Add(element Element)
	Size() int
	Reset()
	Edge() (min, max Element)
}

type Series struct {
	elements []Element
	index    int
	dim      int
	min      Element
	max      Element
}

func NewSeries(dim int) *Series {
	return &Series{dim: dim, elements: make([]Element, 0), min: *NewPoint("min", make([]float64, dim)...), max: *NewPoint("max", make([]float64, dim)...)}
}

func (s *Series) Reset() {
	s.index = 0
}

func (s *Series) Size() int {
	return len(s.elements)
}

func (s *Series) Add(element Element) {
	if element.Dim() != s.dim {
		log.Fatalf("cannot add to Series of dimensionality %d element of dimension %d: %v", s.dim, element.Dim(), element)
	}
	s.elements = append(s.elements, element)

	for i, c := range element.Coords() {
		if c < s.min.Coords()[i] {
			s.min.Coords()[i] = c
		}
		if c > s.max.Coords()[i] {
			s.max.Coords()[i] = c
		}
	}
}

func (s *Series) Edge() (min, max Element) {
	return s.min, s.max
}

func (s *Series) Next() (element Element, hasNext bool) {
	l := len(s.elements)

	if l > s.index {
		oldIndex := s.index
		s.index++
		return s.elements[oldIndex], l > s.index
	}
	return nil, false
}
