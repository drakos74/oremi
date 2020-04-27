package style

import (
	"github/drakos74/oremi/internal/gui"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Checkbox struct {
	gui.RawItem
	checkbox *widget.CheckBox
}

func NewCheckBox() *Checkbox {
	return &Checkbox{
		*gui.NewRawItem(),
		new(widget.CheckBox),
	}
}

func (c *Checkbox) Draw(gtx *layout.Context, th *material.Theme) error {
	th.CheckBox("Checkbox").Layout(gtx, c.checkbox)
	return nil
}
