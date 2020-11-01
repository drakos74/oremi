package gui

import (
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/unit"
	"gioui.org/widget"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

type View struct {
	itemsList
	*layout.List
	height float32
	width  float32
}

func NewView(orientation layout.Axis) *View {
	return &View{
		itemsList: itemsList{
			items: make(map[uint32]item),
			index: make([]uint32, 0),
		},
		List: &layout.List{
			Axis: orientation,
		},
	}
}

func (v *View) Draw(gtx layout.Context, th *material.Theme) (layout.Dimensions, error) {
	//TODO : do recover
	children := make([]layout.Widget, len(v.items))
	// generate a widget for each item
	for i := 0; i < len(v.items); i++ {
		children[i] = func(j int) layout.Widget {
			return v.get(j).draw(gtx, th)
		}(i)
	}
	// draw the widgets in a list.
	border := widget.Border{Color: color.RGBA{A: 0xff}, Width: unit.Px(1)}
	return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return v.List.Layout(gtx, len(children), func(gtx layout.Context, i int) layout.Dimensions {
			return children[i](gtx)
		})
	}), nil
}

func (v *View) Event(e *pointer.Event) (redraw bool, err error) {
	for i := 0; i < len(v.items); i++ {
		if v.get(i).event(e) {
			redraw = true
		}
	}
	return redraw, nil
}

func (v *View) WithMaxHeight(height float32) *View {
	v.height = height
	return v
}

func (v *View) WithMaxWidth(width float32) *View {
	v.width = width
	return v
}
