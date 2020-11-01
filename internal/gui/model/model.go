package model

import (
	"github.com/drakos74/oremi/internal/gui/style"
	"github.com/drakos74/oremi/label"

	"github.com/drakos74/oremi/internal/data/model"
	"github.com/drakos74/oremi/internal/math"

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
	Labels() []label.Label
	Style() style.Properties
}

type Series struct {
	model.Collection
	style       style.Properties
	aggregation int
}

func (s Series) Style() style.Properties {
	return s.style
}

func NewSeries(collection model.Collection, style style.Properties, aggregation int) Collection {
	return Series{
		Collection:  collection,
		style:       style,
		aggregation: aggregation,
	}
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
	// aggregate the number of points with an average
	n := s.aggregation
	vv := make([]model.Vector, n)

	var p model.Vector
	for i := 0; i < n; i++ {
		if p, ok, next = s.Collection.Next(); ok {
			vv[i] = p
		}
	}

	// TODO : make a math operation for this
	var x float64
	var y float64
	var nn float64
	for _, v := range vv {
		if len(v.Coords) > 0 {
			x += v.Coords[0]
			y += v.Coords[1]
			nn++
		}
	}

	return &LabeledPoint{
		// TODO : make the coordinate choice connected to the labels and the graph options in general
		Point: f32.Point{
			X: math.Float32(x / nn),
			Y: math.Float32(y / nn),
		},
		Label: p.Label,
	}, ok, next
}

func (s Series) Reset() {
	s.Collection.Reset()
}
