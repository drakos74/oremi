package canvas

import (
	"github/drakos74/oremi/internal/gui"
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

// Container represents a ui scene
type Container struct {
	RawCompoundElement
	gui.InteractiveElement
	rect *f32.Rectangle
}

func (c *Container) Layout(ops *op.Ops) layout.Dimensions {
	return layout.Dimensions{
		Size: image.Point{
			X: int(c.rect.Max.X + 2*gui.Inset),
			Y: int(c.rect.Max.Y + 2*gui.Inset),
		},
		Baseline: 0,
	}
}

// Draw propagates the draw call to all the scene chldren
func (c *Container) Draw(gtx *layout.Context, th *material.Theme) error {
	gtx.Dimensions = c.Layout(gtx.Ops)
	_, err := c.Elements(DrawAction(gtx, th))
	return err
}

// Event propagates a pointer event to all the scene chldren
func (c *Container) Event(e *pointer.Event) (bool, error) {
	ok, err := c.InteractiveElement.Event(e)
	if err != nil {
		return false, err
	}
	if ok {
		return c.Elements(EventAction(e))
	}
	return false, nil
}

// NewContainer creates a new scene
func NewContainer(rect *f32.Rectangle) *Container {
	return &Container{
		*NewCompoundElement(),
		*gui.NewInteractiveElement(*rect),
		rect,
	}
}
