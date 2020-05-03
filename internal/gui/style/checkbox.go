package style

import (
	"github/drakos74/oremi/internal/gui"
	"github/drakos74/oremi/internal/gui/canvas"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// CheckboxControl is a checkbox to be used as a canvas.Control
type CheckboxControl struct {
	gui.RawItem
	label    string
	checkbox *widget.CheckBox
	active   bool
	trigger  chan canvas.Event
	ack      chan canvas.Event
}

func NewCheckBox(label string, active bool) *CheckboxControl {
	cb := &CheckboxControl{
		*gui.NewRawItem(),
		label,
		new(widget.CheckBox),
		active,
		make(chan canvas.Event),
		make(chan canvas.Event),
	}
	cb.checkbox.SetChecked(active)
	return cb
}

func (c *CheckboxControl) Draw(gtx *layout.Context, th *material.Theme) error {
	th.CheckBox(c.label).Layout(gtx, c.checkbox)
	active := c.active
	c.active = c.checkbox.Checked(gtx)
	if c.active != active {
		c.trigger <- canvas.Trigger
	}
	return nil
}

func (c *CheckboxControl) IsActive() bool {
	return c.active
}

func (c *CheckboxControl) Trigger() canvas.EventReceiver {
	return c.trigger
}

func (c *CheckboxControl) Ack() canvas.EventEmitter {
	return c.ack
}
