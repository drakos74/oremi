package graph

import (
	"fmt"
	"github/drakos74/oremi/internal/gui"
	"github/drakos74/oremi/internal/gui/canvas"
	"github/drakos74/oremi/internal/gui/canvas/math"
	"github/drakos74/oremi/internal/gui/style"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

// Axis is an axis element for graphs
type Axis struct {
	canvas.RawElement
	// TODO : axis should not be a container ... as container becomes now a style element
	canvas.Container
	label style.Label
}

// NewAxisX creates a new x axis
func NewAxisX(label string, start f32.Point, length float32, delim int, calc math.Mapper) *Axis {
	rect := &f32.Rectangle{
		Min: start,
		Max: f32.Point{
			X: start.X + length,
			Y: start.Y + 1,
		},
	}
	axis := &Axis{
		*canvas.NewRawElement(),
		*canvas.NewContainer(rect),
		style.NewLabel(f32.Point{
			X: rect.Max.X,
			Y: rect.Max.Y + 20,
		}, label),
	}
	for i := 0; i <= delim; i++ {
		d := float32(i) / float32(delim)
		rect := axis.Rect()
		axis.Add(NewDelimiterX(
			f32.Point{
				X: rect.Min.X + (rect.Max.X-rect.Min.X)*d,
				Y: start.Y,
			},
			calc.DeScaleX()))
	}
	return axis
}

// NewAxisY creates a new y axis
func NewAxisY(label string, start f32.Point, length float32, delim int, calc math.Mapper) *Axis {
	rect := &f32.Rectangle{
		Min: start,
		Max: f32.Point{
			X: start.X + 1,
			Y: start.Y + length,
		},
	}
	axis := &Axis{
		*canvas.NewRawElement(),
		*canvas.NewContainer(rect),
		style.NewLabel(f32.Point{
			X: rect.Min.X - 20,
			Y: rect.Min.Y - 2,
		}, label),
	}
	for i := 0; i <= delim; i++ {
		d := float32(i) / float32(delim)
		rect := axis.Rect()
		axis.Add(NewDelimiterY(
			f32.Point{
				X: start.X,
				Y: rect.Min.Y + (rect.Max.Y-rect.Min.Y)*d,
			},
			calc.DeScaleY()))
	}
	return axis
}

// Draw draws the axis
func (a *Axis) Draw(gtx *layout.Context, th *material.Theme) error {
	paint.PaintOp{Rect: a.Container.Rect()}.Add(gtx.Ops)
	_, err := a.Elements(canvas.DrawAction(gtx, th))
	if err != nil {
		return err
	}
	return a.label.Draw(gtx, th)
}

// Delimiter is an axis child element representing a value on the respective axis
type Delimiter struct {
	canvas.RawElement
	gui.InteractiveElement
	rect      f32.Rectangle
	label     style.Label
	transform func() float32
}

// NewDelimiterX creates a new delimiter for an x axis
func NewDelimiterX(p f32.Point, transform math.Transform) *Delimiter {
	rect := f32.Rectangle{
		Min: f32.Point{
			X: p.X,
			Y: p.Y - 10,
		},
		Max: f32.Point{
			X: p.X + 1,
			Y: p.Y + 10,
		},
	}
	return &Delimiter{
		*canvas.NewRawElement(),
		*gui.NewInteractiveElement(rect),
		rect,
		style.NewLabel(p, ""),
		func() float32 {
			return transform(p.X)
		},
	}
}

// NewDelimiterY creates a new delimiter for an x axis
func NewDelimiterY(p f32.Point, transform math.Transform) *Delimiter {
	rect := f32.Rectangle{
		Min: f32.Point{
			X: p.X - 10,
			Y: p.Y,
		},
		Max: f32.Point{
			X: p.X + 10,
			Y: p.Y + 1,
		},
	}
	return &Delimiter{
		*canvas.NewRawElement(),
		*gui.NewInteractiveElement(rect),
		rect,
		style.NewLabel(p, ""),
		func() float32 {
			return transform(p.Y)
		},
	}
}

// Draw draws the delimiter
func (m *Delimiter) Draw(gtx *layout.Context, th *material.Theme) error {
	r := m.rect
	if m.IsActive() {
		m.label.Text(fmt.Sprintf("%v", m.transform()))
		err := m.label.Draw(gtx, th)
		if err != nil {
			return err
		}
	}
	paint.PaintOp{Rect: r}.Add(gtx.Ops)
	return nil
}
