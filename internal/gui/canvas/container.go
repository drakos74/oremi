package canvas

import (
	"image"

	"gioui.org/unit"
	"github.com/drakos74/oremi/internal/gui"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

// Container represents a ui scene
type Container struct {
	*gui.InteractiveElement
	CompoundElement
}

// Draw propagates the draw call to all the scene children.
func (c *Container) Draw(gtx layout.Context, th *material.Theme) (layout.Dimensions, error) {
	//layout.Dimensions = c.Layout(gtx.Ops)
	p, active := c.InteractiveElement.Pointer()
	var err error
	if active {
		_, err = c.Elements(gtx, EventAction(p), DrawAction(gtx, th))
	} else {
		_, err = c.Elements(gtx, DrawAction(gtx, th))
	}
	return layout.Dimensions{
		Size: image.Point{
			X: gtx.Px(unit.Px(c.InteractiveElement.Area.Rect().Max.X - c.InteractiveElement.Area.Rect().Min.X)),
			Y: gtx.Px(unit.Px(c.InteractiveElement.Area.Rect().Max.Y - c.InteractiveElement.Area.Rect().Min.Y)),
		},
	}, err
}

// NewContainer creates a new scene
func NewContainer(rect *f32.Rectangle) *Container {
	return &Container{
		gui.NewInteractiveElement(rect),
		NewCompoundElement(),
	}
}
