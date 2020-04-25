package model

import (
	"github/drakos74/oremi/internal/data/model"

	"gioui.org/f32"
)

type LabeledPoint struct {
	f32.Point
	Label []string
}

// TODO : embed the data model collection into the interface
type Collection interface {
	Bounds() *f32.Rectangle
	Next() (point *LabeledPoint, ok, next bool)
	Size() int
	Reset()
	Labels() []string
}

type Series struct {
	model.Collection
}

func NewSeries(collection model.Collection) Collection {
	return Series{collection}
}

func (s Series) Bounds() *f32.Rectangle {
	min, max := s.Edge()
	return &f32.Rectangle{
		Min: f32.Point{
			X: float32(min.Coords[0]),
			Y: float32(min.Coords[1]),
		},
		Max: f32.Point{
			X: float32(max.Coords[0]),
			Y: float32(max.Coords[1]),
		},
	}
}

func (s Series) Next() (point *LabeledPoint, ok, next bool) {
	if p, ok, next := s.Collection.Next(); ok {
		return &LabeledPoint{
			// TODO : make the coordinate choice connected to the labels and the graph options in general
			Point: f32.Point{
				X: float32(p.Coords[0]),
				Y: float32(p.Coords[1]),
			},
			Label: p.Label,
		}, true, next
	}
	return nil, false, false
}

func (s Series) Reset() {
	s.Collection.Reset()
}
