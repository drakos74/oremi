package style

import (
	"github/drakos74/oremi/internal/gui"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type CheckboxControl struct {
	gui.RawItem
	label    string
	checkbox *widget.CheckBox
}

func NewCheckBox(label string) *CheckboxControl {
	cb := &CheckboxControl{
		*gui.NewRawItem(),
		label,
		new(widget.CheckBox),
	}
	cb.checkbox.SetChecked(true)
	return cb
}

func (c *CheckboxControl) Draw(gtx *layout.Context, th *material.Theme) error {
	th.CheckBox(c.label).Layout(gtx, c.checkbox)
	return nil
}

func (c *CheckboxControl) IsActive(gtx *layout.Context) bool {
	return c.checkbox.Checked(gtx)
}
