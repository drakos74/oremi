package style

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Input struct {
	editor *widget.Editor
}

func NewInput() *Input {
	return &Input{&widget.Editor{
		SingleLine: true,
		Submit:     true,
	}}
}

func (i Input) Draw(gtx *layout.Context, th *material.Theme) error {
	e := th.Editor("Hint")
	e.Font.Style = text.Italic
	e.Layout(gtx, i.editor)
	for _, e := range i.editor.Events(gtx) {
		if _, ok := e.(widget.SubmitEvent); ok {
			i.editor.SetText("")
		}
	}
	return nil
}
