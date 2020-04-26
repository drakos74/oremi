package gui

import (
	"log"

	"github.com/google/uuid"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type items map[uint32]item

type item struct {
	draw  func(gtx *layout.Context, th *material.Theme) func()
	event func(e *pointer.Event) bool
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
	eventAction := func(e *pointer.Event) bool {
		return false
	}
	if ev, ok := v.(InteractiveItem); ok {
		eventAction = func(e *pointer.Event) bool {
			a, err := ev.Event(e)
			if err != nil {
				log.Fatalf("error encountered during Event Propagation: %v", err)
			}
			return a
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

// DrawItem represents an elements that can be drawn on the screen
type DrawItem interface {
	Draw(gtx *layout.Context, th *material.Theme) error
}

// InteractiveItem represents an interactive UI element
type InteractiveItem interface {
	Event(e *pointer.Event) (bool, error)
	Activate() error
	Reset() error
	IsActive() bool
	Rect() f32.Rectangle
}

// InteractiveElement is the base implementation of a dynamic element
type InteractiveElement struct {
	active bool
	halo   float32
	bounds f32.Rectangle
}

// Event propagates the scene event to the element
func (s *InteractiveElement) Event(e *pointer.Event) (bool, error) {
	cState := s.active
	if !checkRect(s.bounds, e.Position, s.halo) {
		err := s.Reset()
		return cState == false, err
	} else {
		err := s.Activate()
		return cState && err == nil, err
	}
}

// Activate triggers the activation of a dynamic element
func (s *InteractiveElement) Activate() error {
	s.active = true
	return nil
}

// Reset resets the state of an dynamic element
func (s *InteractiveElement) Reset() error {
	s.active = false
	return nil
}

// IsActive returns the activation status of a dynamic element
func (s *InteractiveElement) IsActive() bool {
	return s.active
}

func (s *InteractiveElement) Rect() f32.Rectangle {
	return s.bounds
}

// NewInteractiveElement creates a new dynamic element
func NewInteractiveElement(rect f32.Rectangle) *InteractiveElement {
	return &InteractiveElement{
		halo:   4,
		bounds: rect,
	}
}

func checkRect(rect f32.Rectangle, p f32.Point, s float32) bool {
	r := f32.Rectangle{
		Min: f32.Point{X: p.X - s, Y: p.Y - s},
		Max: f32.Point{X: p.X + s, Y: p.Y + s},
	}
	return !rect.Intersect(r).Empty()
}
