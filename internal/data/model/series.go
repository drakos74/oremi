package model

import (
	"fmt"
	"log"
	"math"

	"github.com/drakos74/oremi/label"
)

// Series is a collection of vectors
type Series struct {
	vectors []Vector
	index   int
	dim     int
	min     Vector
	max     Vector
	labels  []label.Label
	events  chan Event
}

// NewSeries creates a new series of the specified dimension
func NewSeries(labels ...label.Label) *Series {
	dim := len(labels)
	min := make([]float64, dim)
	for i := range min {
		min[i] = math.MaxFloat64
	}
	return &Series{
		dim:     dim,
		vectors: make([]Vector, 0),
		min:     NewVector([]string{"min"}, min...),
		max:     NewVector([]string{"max"}, make([]float64, dim)...),
		labels:  labels,
		events:  make(chan Event, 100),
	}
}

func (s *Series) Events() <-chan Event {
	return s.events
}

// Reset resets the iterator to the start of the collection
func (s *Series) Reset() {
	s.index = 0
}

// Next returns the next vector in the series
func (s *Series) Next() (vector Vector, ok, hasNext bool) {
	l := len(s.vectors)

	if l > s.index {
		oldIndex := s.index
		s.index++
		return s.vectors[oldIndex], true, l > s.index
	}
	return Vector{}, false, false
}

// Size returns the size of the series
func (s *Series) Size() int {
	return len(s.vectors)
}

// Add adds a vector to the series
// the call should fail if the vectors dimensions are not the same as the ones of the defined series
func (s *Series) Add(vector Vector) {
	if vector.Dim() != s.dim {
		log.Fatalf("cannot add to Series of dimensionality %d vector of dimension %d: %v", s.dim, vector.Dim(), vector)
	}
	s.vectors = append(s.vectors, vector)

	for i, c := range vector.Coords {
		if c < s.min.Coords[i] {
			s.min.Coords[i] = c
		}
		if c > s.max.Coords[i] {
			s.max.Coords[i] = c
		}
	}
	s.events <- Event{
		T: Added,
		A: true,
		S: fmt.Sprintf("%+v", vector),
	}
}

// Edge returns the edge values of the series
// this is useful for quick comparisons of collections of data, as well as drawing and scaling
func (s *Series) Edge() (min, max Vector) {
	return s.min, s.max
}

// Labels returns the labels of the series.
func (s *Series) Labels() []label.Label {
	return s.labels
}

// String returns the data arrays of the series.
func (s Series) String() string {
	return fmt.Sprintf("%v", s.vectors)
}
