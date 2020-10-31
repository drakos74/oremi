package gui

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"

	"gioui.org/f32"
	"gioui.org/op/paint"

	"gioui.org/unit"

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

	children := make([]layout.FlexChild, len(v.items))

	for i := 0; i < len(v.items); i++ {
		children[i] = layout.Rigid(func(j int) layout.Widget {
			item := v.get(j)
			d := layout.Dimensions{}
			//if v, ok := item.(View); ok {
			//	d = layout.Dimensions{
			//		Size: image.Point{
			//			X: gtx.Px(unit.Dp(v.width)),
			//			Y: gtx.Px(unit.Dp(v.height)),
			//		},
			//	}
			//}
			return func(gtx layout.Context) layout.Dimensions {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, item.draw(gtx, th))
				return d
			}
		}(i))
	}

	layout.Flex{Alignment: layout.Start}.Layout(gtx, children...)

	list := &layout.List{
		Axis: layout.Vertical,
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			paint.ColorOp{Color: color.RGBA{127, 0, 0, 255}}.Add(gtx.Ops)
			paint.PaintOp{Rect: f32.Rect(0, 0, 100, 100)}.Add(gtx.Ops)
			paint.ColorOp{Color: color.RGBA{0, 127, 0, 255}}.Add(gtx.Ops)
			paint.PaintOp{Rect: f32.Rect(50, 50, 150, 150)}.Add(gtx.Ops)

			return layout.Dimensions{
				//Size: image.Point{
				//	X: 150,
				//	Y: 150,
				//},
			}
		},
	}

	list.Layout(gtx, len(widgets), func(gtx layout.Context, i int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[i])
	})

	//v.Layout(gtx, len(v.items), func(gtx layout.Context, i int) layout.Dimensions {
	//	return layout.UniformInset(unit.Dp(0)).Layout(gtx, v.get(i).draw(gtx, th))
	//})
	return layout.Dimensions{
		Size: image.Point{
			X: gtx.Px(unit.Dp(v.width)),
			Y: gtx.Px(unit.Dp(v.height)),
		},
	}, nil
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
