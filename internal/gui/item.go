package gui

import (
	"image"
	"log"
	"reflect"

	"gioui.org/op"

	"gioui.org/unit"

	"github.com/google/uuid"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type items map[uint32]item

type item struct {
	kind  string
	size  f32.Point
	draw  func(gtx layout.Context, th *material.Theme) layout.Widget
	event func(e *pointer.Event) bool
}

// Add accepts a generic object that must have one of the following properties :
// - implement
func (items items) AddItem(v interface{}) (uint32, func()) {
	item := newItem(v)
	//item.size = size
	id := uuid.New().ID()
	items[id] = item
	return id, func() {
		delete(items, id)
	}
}

// TODO : measure interface perf overhead usage
func newItem(v interface{}) item {
	item := item{
		kind: reflect.TypeOf(v).String(),
	}
	// apply draw action if applicable
	drawAction := func(gtx layout.Context, th *material.Theme) layout.Widget {
		return func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{}
		}
	}
	if d, ok := v.(DrawItem); ok {
		drawAction = func(gtx layout.Context, th *material.Theme) layout.Widget {
			return func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					d, err := d.Draw(gtx, th)
					if err != nil {
						log.Fatalf("error encountered during Draw: %v", err)
					}
					return d
				})
			}
		}
	}
	item.draw = drawAction

	// apply event listener if applicable
	eventAction := func(e *pointer.Event) bool {
		return false
	}
	if ev, ok := v.(InteractiveItem); ok {
		item.draw = func(gtx layout.Context, th *material.Theme) layout.Widget {
			err := ev.Enable(gtx)
			if err != nil {
				log.Fatalf("could not enable interactive element: %v: %v", ev, err)
			}
			return drawAction(gtx, th)
		}
		eventAction = func(e *pointer.Event) bool {
			err := ev.Event(e)
			if err != nil {
				log.Fatalf("error encountered during Event Propagation: %v", err)
			}
			return ok
		}
	}
	if view, ok := v.(*View); ok {
		eventAction = func(e *pointer.Event) bool {
			ok, err := view.Event(e)
			if err != nil {
				log.Fatalf("error encountered during Event Propagation: %v", err)
			}
			return ok
		}
	}
	item.event = eventAction

	return item
}

type itemsList struct {
	items
	index []uint32
}

func (list *itemsList) Add(vv ...interface{}) func() {
	rmvs := make([]func(), len(vv))
	for i, v := range vv {
		id, rmv := list.items.AddItem(v)
		list.index = append(list.index, id)
		rmvs[i] = func() {
			rmv()
			rk := -1
			for i, id := range list.index {
				if _, exists := list.items[id]; !exists {
					rk = i
					break
				}
			}
			copy(list.index[rk:], list.index[rk+1:])    // Shift a[i+1:] left one index.
			list.index = list.index[:len(list.index)-1] // Truncate slice.
		}
	}
	return func() {
		for _, rm := range rmvs {
			rm()
		}
	}
}

func (list itemsList) get(i int) item {
	return list.items[list.index[i]]
}

// Item is the main abstraction for any object living within the canvas
type Item interface {
	ID() uint32
}

// RawItem is the base implementation for an Item
type RawItem struct {
	id uint32
}

// ID returns the id of the raw element
func (item *RawItem) ID() uint32 {
	return item.id
}

// NewRawItem creates a new raw element
func NewRawItem() *RawItem {
	return &RawItem{
		id: uuid.New().ID(),
	}
}

// Area is the main abstraction for any object taking up or actively tracking some space on the canvas
type Area interface {
	Rect() f32.Rectangle
	Expand(halo int) image.Rectangle
	Size(inset int) image.Point
}

// RawItem is the base implementation for an Item
type RawArea struct {
	rect *f32.Rectangle
}

func (area *RawArea) Expand(halo int) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(area.rect.Min.X) - halo,
			Y: int(area.rect.Min.Y) - halo,
		},
		Max: image.Point{
			X: int(area.rect.Max.X) + halo,
			Y: int(area.rect.Max.Y) + halo,
		},
	}
}

func (area *RawArea) Size(inset int) image.Point {
	return image.Point{
		X: int(area.rect.Max.X) + 2*inset,
		Y: int(area.rect.Max.Y) + 2*inset,
	}
}

const ff = 1

// ID returns the id of the raw element
func (area *RawArea) Rect() f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{
			X: ff * area.rect.Min.X,
			Y: ff * area.rect.Min.Y,
		},
		Max: f32.Point{
			X: ff * area.rect.Max.X,
			Y: ff * area.rect.Max.Y,
		},
	}
}

// Rect creates a new raw element
func Rect(rect *f32.Rectangle) *RawArea {
	return &RawArea{
		rect: rect,
	}
}

// DrawItem represents an elements that can be drawn on the screen
type DrawItem interface {
	Draw(gtx layout.Context, th *material.Theme) (layout.Dimensions, error)
}

// InteractiveItem represents an interactive UI element
type InteractiveItem interface {
	Event(e *pointer.Event) error
	Enable(gtx layout.Context) error
	Pointer() (f32.Point, bool)
}

// InteractiveElement is the base implementation of a dynamic element
type InteractiveElement struct {
	Item
	Area
	halo    int
	active  bool
	pointer f32.Point
}

// NewInteractiveElement creates a new dynamic element
func NewInteractiveElement(rect *f32.Rectangle) *InteractiveElement {
	return &InteractiveElement{
		Item: NewRawItem(),
		Area: Rect(rect),
		halo: 4,
	}
}

// Enable adds the event handler for an interactive element
func (item *InteractiveElement) Enable(gtx layout.Context) error {
	// TODO : fix the event handling to enable active check
	stack := op.Push(gtx.Ops)
	rect := item.Expand(item.halo)
	pointer.Rect(rect).Add(gtx.Ops)
	pointer.InputOp{Tag: item.ID(), Types: pointer.Enter}.Add(gtx.Ops)
	pointer.InputOp{Tag: item.ID(), Types: pointer.Leave}.Add(gtx.Ops)
	stack.Pop()

	events := gtx.Events(item.ID())
	if len(events) > 0 {
		event := events[0]
		if ev, ok := event.(pointer.Event); ok {
			switch ev.Type {
			case pointer.Enter:
				item.active = true
			case pointer.Leave:
				item.active = false
			}
		}
	}
	return nil
}

// Event propagates the scene event to the element
func (item *InteractiveElement) Event(e *pointer.Event) error {
	item.pointer = e.Position
	return nil
}

func (item *InteractiveElement) Pointer() (f32.Point, bool) {
	return item.pointer, item.active
}
