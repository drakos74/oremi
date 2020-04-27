package gui

import (
	"image"
	"log"

	"gioui.org/op"

	"github.com/google/uuid"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type items map[uint32]item

type item struct {
	draw  func(gtx *layout.Context, th *material.Theme) func()
	event func(gtx *layout.Context, e *pointer.Event) bool
}

// AddItem accepts a generic object that must have one of the following properties :
// - implement
func (items items) AddItem(v interface{}) func() {
	item := newItem(v)
	id := uuid.New().ID()
	items[id] = item
	return func() {
		delete(items, id)
	}
}

func (items items) idx() []uint32 {
	idx := make([]uint32, len(items))
	i := 0
	for itemIndex := range items {
		idx[i] = itemIndex
		i++
	}
	return idx
}

// TODO : measure interface perf overhead usage
func newItem(v interface{}) item {
	item := item{}
	// apply draw action if applicable
	drawAction := func(gtx *layout.Context, th *material.Theme) func() {
		return func() {
			// void
		}
	}
	if d, ok := v.(DrawItem); ok {
		drawAction = func(gtx *layout.Context, th *material.Theme) func() {
			return func() {
				err := d.Draw(gtx, th)
				if err != nil {
					log.Fatalf("error encountered during Draw: %v", err)
				}
			}
		}
	}
	item.draw = drawAction

	// apply event listener if applicable
	eventAction := func(gtx *layout.Context, e *pointer.Event) bool {
		return false
	}
	if ev, ok := v.(InteractiveItem); ok {
		item.draw = func(gtx *layout.Context, th *material.Theme) func() {
			err := ev.Enable(gtx, th)
			if err != nil {
				log.Fatalf("could not enable interactive element: %v: %v", ev, err)
			}
			return drawAction(gtx, th)
		}
		eventAction = func(gtx *layout.Context, e *pointer.Event) bool {
			_, ok, err := ev.Event(gtx, e)
			if err != nil {
				log.Fatalf("error encountered during Event Propagation: %v", err)
			}
			return ok
		}
	}
	item.event = eventAction
	return item
}

type list struct {
	items
	index []uint32
}

func (list *list) AddItem(v interface{}) func() {
	rmv := list.items.AddItem(v)
	list.index = list.items.idx()
	return func() {
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

func (list list) idx() []uint32 {
	return list.index
}

func (list list) get(i int) item {
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

// ID returns the id of the raw element
func (area *RawArea) Rect() f32.Rectangle {
	return *area.rect
}

// Rect creates a new raw element
func Rect(rect *f32.Rectangle) *RawArea {
	return &RawArea{
		rect: rect,
	}
}

// DrawItem represents an elements that can be drawn on the screen
type DrawItem interface {
	Draw(gtx *layout.Context, th *material.Theme) error
}

// TODO : cleanup interface
// InteractiveItem represents an interactive UI element
type InteractiveItem interface {
	Event(gtx *layout.Context, e *pointer.Event) (f32.Point, bool, error)
	Enable(gtx *layout.Context, th *material.Theme) error
	Activate() error
	Reset() error
	IsActive() bool
}

// InteractiveElement is the base implementation of a dynamic element
type InteractiveElement struct {
	Item
	Area
	active bool
	halo   int
	offset f32.Point
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
func (item *InteractiveElement) Enable(gtx *layout.Context, th *material.Theme) error {
	var stack op.StackOp
	rect := item.Expand(item.halo)
	stack.Push(gtx.Ops)
	pointer.Rect(rect).Add(gtx.Ops)
	pointer.InputOp{Key: item.ID()}.Add(gtx.Ops)
	stack.Pop()
	return nil
}

// Event propagates the scene event to the element
func (item *InteractiveElement) Event(gtx *layout.Context, e *pointer.Event) (f32.Point, bool, error) {
	events := gtx.Events(item.ID())
	if len(events) > 0 {
		event := events[0]
		if ev, ok := event.(pointer.Event); ok {
			return ev.Position, true, item.Activate()
		}
		return f32.Point{}, false, nil
	}
	if item.IsActive() && len(events) == 0 {
		return f32.Point{}, true, item.Reset()
	}
	return f32.Point{}, false, nil
}

// Activate triggers the activation of a dynamic element
func (item *InteractiveElement) Activate() error {
	item.active = true
	return nil
}

// Reset resets the state of an dynamic element
func (item *InteractiveElement) Reset() error {
	item.active = false
	return nil
}

// IsActive returns the activation status of a dynamic element
func (item *InteractiveElement) IsActive() bool {
	return item.active
}
