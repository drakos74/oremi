package style

import (
	"github/drakos74/oremi/internal/gui/canvas"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Checkbox struct {
	canvas.RawElement
	checkbox *widget.CheckBox
}

func NewCheckBox() *Checkbox {
	return &Checkbox{
		*canvas.NewRawElement(),
		new(widget.CheckBox),
	}
}

func (c *Checkbox) Draw(gtx *layout.Context, th *material.Theme) error {
	th.CheckBox("Checkbox").Layout(gtx, c.checkbox)
	return nil
}
