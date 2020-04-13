package entity

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"oremi/internal/canvas"
)

type Scene struct {
	canvas.RawCompoundElement
	canvas.RawDynamicElement
	rect *f32.Rectangle
}

func (s *Scene) Draw(gtx *layout.Context, th *material.Theme) error {
	_, err := s.Elements(canvas.DrawAction(gtx, th))
	return err
}

func (s *Scene) Event(e *pointer.Event) (bool, error) {
	ok, err := s.RawDynamicElement.Event(e)
	if err != nil {
		return false, err
	}
	if ok {
		return s.Elements(canvas.EventAction(e))
	}
	return false, nil
}

func NewScene(rect *f32.Rectangle) *Scene {
	return &Scene{
		*canvas.NewCompoundElement(),
		*canvas.NewDynamicElement(*rect),
		rect,
	}
}
