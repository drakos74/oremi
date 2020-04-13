package entity

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"oremi/internal/canvas"
)

type Point struct {
	canvas.RawElement
	canvas.RawDynamicElement
	w     float32
	c     f32.Point
	rect  *f32.Rectangle
	label string
}

func NewPoint(label string, center f32.Point) *Point {
	var w float32 = 2
	rect := calculateRect(center, w)
	p := &Point{
		*canvas.NewRawElement(),
		*canvas.NewDynamicElement(rect),
		w,
		center,
		&rect,
		label,
	}
	return p
}

func calculateRect(center f32.Point, w float32) f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{X: center.X - w, Y: center.Y - w},
		Max: f32.Point{X: center.X + w, Y: center.Y + w},
	}
}

func (p *Point) Draw(gtx *layout.Context, th *material.Theme) error {
	r := *p.rect
	if p.IsActive() {
		r = calculateRect(p.c, 2*p.w)
		println(fmt.Sprintf("%s", p.label))
		// TODO : fix show label on hover and remove println
		th.Label(unit.Px(50), p.label).Layout(gtx)
	}
	paint.PaintOp{Rect: r}.Add(gtx.Ops)

	return nil
}
