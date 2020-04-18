package entity

import (
	"github/drakos74/oremi/internal/gui/canvas"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

// Container represents a ui scene
type Container struct {
	canvas.RawCompoundElement
	canvas.RawDynamicElement
	rect *f32.Rectangle
}

// Draw propagates the draw call to all the scene chldren
func (s *Container) Draw(gtx *layout.Context, th *material.Theme) error {
	_, err := s.Elements(canvas.DrawAction(gtx, th))
	return err
}

// Event propagates a pointer event to all the scene chldren
func (s *Container) Event(e *pointer.Event) (bool, error) {
	ok, err := s.RawDynamicElement.Event(e)
	if err != nil {
		return false, err
	}
	if ok {
		return s.Elements(canvas.EventAction(e))
	}
	return false, nil
}

// NewContainer creates a new scene
func NewContainer(rect *f32.Rectangle) *Container {
	return &Container{
		*canvas.NewCompoundElement(),
		*canvas.NewDynamicElement(*rect),
		rect,
	}
}
