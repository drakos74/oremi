package graph

import (
	"github/drakos74/oremi/internal/gui"
	"github/drakos74/oremi/internal/gui/canvas"
	"github/drakos74/oremi/internal/gui/style"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

// Point is a point element
type Point struct {
	canvas.RawElement
	gui.InteractiveElement
	w     float32
	c     f32.Point
	rect  *f32.Rectangle
	label style.Label
}

// NewPoint creates a new point
func NewPoint(label string, center f32.Point) *Point {
	var w float32 = 2
	rect := calculateRect(center, w)
	p := &Point{
		*canvas.NewRawElement(),
		*gui.NewInteractiveElement(rect),
		w,
		center,
		&rect,
		style.NewLabel(center, label),
	}
	return p
}

func calculateRect(center f32.Point, w float32) f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{X: center.X - w, Y: center.Y - w},
		Max: f32.Point{X: center.X + w, Y: center.Y + w},
	}
}

// Draw draws the point on the canvas
func (p *Point) Draw(gtx *layout.Context, th *material.Theme) error {
	r := *p.rect
	if p.IsActive() {
		r = calculateRect(p.c, 2*p.w)
		err := p.label.Draw(gtx, th)
		if err != nil {
			return err
		}
	}
	paint.PaintOp{Rect: r}.Add(gtx.Ops)
	return nil
}
