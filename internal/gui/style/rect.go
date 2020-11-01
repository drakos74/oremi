package style

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

type Rect struct {
	rect *f32.Rectangle
}

func NewRect(rect *f32.Rectangle) *Rect {
	return &Rect{rect: rect}
}

func (r Rect) Draw(gtx layout.Context, th *material.Theme) (layout.Dimensions, error) {
	paint.PaintOp{Rect: *r.rect}.Add(gtx.Ops)
	return layout.Dimensions{}, nil
}
