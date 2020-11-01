package style

import (
	"image/color"

	"gioui.org/text"
	"gioui.org/unit"
	"github.com/drakos74/oremi/internal/gui/canvas"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Input struct {
	text    string
	editor  *widget.Editor
	trigger chan canvas.Event
}

func NewInput() *Input {
	return &Input{"",
		&widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		make(canvas.Events)}
}

func (i *Input) Draw(gtx layout.Context, th *material.Theme) (layout.Dimensions, error) {
	e := material.Editor(th, i.editor, "Hint")
	e.Font.Style = text.Italic
	border := widget.Border{Color: color.RGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
	for _, e := range i.editor.Events() {
		if _, ok := e.(widget.SubmitEvent); ok {
			i.editor.SetText("")
		}
	}
	if i.text != i.editor.Text() {
		i.text = i.editor.Text()
		i.trigger <- canvas.Event{
			T: canvas.Trigger,
			A: false,
			S: i.text,
		}
	}
	return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
	}), nil
}

func (i *Input) IsActive() bool {
	panic("implement me")
}

func (i *Input) Set(active bool) {
	panic("implement me")
}

func (i *Input) Trigger() canvas.EventReceiver {
	return i.trigger
}

func (i *Input) Ack() canvas.EventEmitter {
	panic("implement me")
}
