package entity

import (
	"fmt"
	"github/drakos74/oremi/internal/gui/canvas"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

type Rect struct {
	canvas.RawElement
	rect *f32.Rectangle
}

func NewRect(x, y, width, height float32) *Rect {
	return &Rect{
		*canvas.NewRawElement(),
		&f32.Rectangle{
			Min: f32.Point{
				X: x,
				Y: y,
			},
			Max: f32.Point{
				X: width,
				Y: height,
			},
		},
	}
}

func (r *Rect) Draw(gtx *layout.Context, th *material.Theme) error {
	paint.PaintOp{Rect: *r.rect}.Add(gtx.Ops)
	return nil
}

func (r *Rect) String() string {
	return fmt.Sprintf("%v", r.rect)
}
