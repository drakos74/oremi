package entity

import (
	"github/drakos74/oremi/internal/canvas"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

type Axis struct {
	canvas.RawElement
	canvas.RawCompoundElement
	rect f32.Rectangle
}

func NewAxisX(start f32.Point, length float32, delim int) *Axis {
	axis := &Axis{
		*canvas.NewRawElement(),
		*canvas.NewCompoundElement(),
		f32.Rectangle{
			Min: start,
			Max: f32.Point{
				X: start.X + length,
				Y: start.Y + 1,
			},
		},
	}
	for i := 0; i <= delim; i++ {
		d := float32(i) / float32(delim)
		axis.Add(NewXDelimiter(f32.Point{
			X: axis.rect.Min.X + (axis.rect.Max.X-axis.rect.Min.X)*d,
			Y: start.Y,
		}))
	}
	return axis
}

func NewAxisY(start f32.Point, length float32, delim int) *Axis {
	axis := &Axis{
		*canvas.NewRawElement(),
		*canvas.NewCompoundElement(),
		f32.Rectangle{
			Min: start,
			Max: f32.Point{
				X: start.X + 1,
				Y: start.Y + length,
			},
		},
	}
	for i := 0; i <= delim; i++ {
		d := float32(i) / float32(delim)
		axis.Add(NewYDelimiter(f32.Point{
			X: start.X,
			Y: axis.rect.Min.Y + (axis.rect.Max.Y-axis.rect.Min.Y)*d,
		}))
	}
	return axis
}

func (a *Axis) Draw(gtx *layout.Context, th *material.Theme) error {
	paint.PaintOp{Rect: a.rect}.Add(gtx.Ops)
	_, err := a.Elements(canvas.DrawAction(gtx, th))
	if err != nil {
		return err
	}
	return nil
}

type Delimiter struct {
	canvas.RawElement
	rect f32.Rectangle
}

func NewXDelimiter(p f32.Point) *Delimiter {
	return &Delimiter{
		*canvas.NewRawElement(),
		f32.Rectangle{
			Min: f32.Point{
				X: p.X,
				Y: p.Y - 10,
			},
			Max: f32.Point{
				X: p.X + 1,
				Y: p.Y + 10,
			},
		},
	}
}

func NewYDelimiter(p f32.Point) *Delimiter {
	return &Delimiter{
		*canvas.NewRawElement(),
		f32.Rectangle{
			Min: f32.Point{
				X: p.X - 10,
				Y: p.Y,
			},
			Max: f32.Point{
				X: p.X + 10,
				Y: p.Y + 1,
			},
		},
	}
}

func (m *Delimiter) Draw(gtx *layout.Context, th *material.Theme) error {
	paint.PaintOp{Rect: m.rect}.Add(gtx.Ops)
	return nil
}
