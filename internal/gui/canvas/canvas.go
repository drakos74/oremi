package canvas

import (
	"github/drakos74/oremi/internal/gui"
	"log"

	"gioui.org/f32"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

// TODO : replace with items abstraction
// CompoundElement represents an element that can have children
type CompoundElement interface {
	Add(elements ...gui.Item)
	Remove(id uint32)
	Elements(apply Action) (bool, error)
}

// Action defines an action to be applied to an Item
type Action func(element gui.Item) (bool, error)

// DrawFunction is a helper method to invoke the Draw method on an elements
var DrawAction = func(gtx *layout.Context, th *material.Theme) func(element gui.Item) (bool, error) {
	return func(element gui.Item) (bool, error) {
		// TODO : avoid reflection by keeping the draw actions in a slice
		if el, ok := element.(gui.DrawItem); ok {
			err := el.Draw(gtx, th)
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}
}

// EventAction is a helper method to invoke the Event method on an elements
var EventAction = func(pointer f32.Point) func(element gui.Item) (bool, error) {
	return func(element gui.Item) (bool, error) {
		var active bool
		// TODO : avoid reflection by keeping the event actions in a slice
		if el, ok := element.(DynamicElement); ok {
			a, err := el.Event(pointer)
			if err != nil {
				return false, err
			}
			if a {
				active = true
			}
		}
		return active, nil
	}
}

// RawCompundElement is the base implementation for a compund element
type RawCompoundElement struct {
	elements map[uint32]gui.Item
}

// AddItem adds a new element to the group
func (s *RawCompoundElement) Add(elements ...gui.Item) {
	//t := reflect.TypeOf(element)
	// TODO : use to make better use of generic actions, without doing the casting
	if elements == nil {
		log.Fatalf("cannot add nil elements")
	}
	for _, element := range elements {
		s.elements[element.ID()] = element
	}
}

// Elements applies the specified action to all child elements
func (s *RawCompoundElement) Elements(apply Action) (bool, error) {
	var d bool
	for _, e := range s.elements {
		done, err := apply(e)
		if err != nil {
			return false, err
		}
		if done {
			d = true
		}
	}
	return d, nil
}

// Remove removes an element by id from the group
func (s *RawCompoundElement) Remove(id uint32) {
	delete(s.elements, id)
}

// Size returns the number of child elements
func (s *RawCompoundElement) Size() int {
	return len(s.elements)
}

// NewCompoundElement creates a new compound element
func NewCompoundElement() *RawCompoundElement {
	return &RawCompoundElement{elements: make(map[uint32]gui.Item)}
}

// DynamicElement represents an interactive UI element
type DynamicElement interface {
	Event(pointer f32.Point) (bool, error)
	Activate() error
	Reset() error
	IsActive() bool
}

// RawDynamicElement is the base implementation of a dynamic element
type RawDynamicElement struct {
	gui.Area
	active bool
	// TODO : leave this to be configured by the child element
	halo float32
}

// Event propagates the scene event to the element
func (s *RawDynamicElement) Event(pointer f32.Point) (bool, error) {
	if !checkRect(s.Rect(), pointer, s.halo) {
		err := s.Reset()
		return false, err
	} else {
		err := s.Activate()
		return err == nil, err
	}
}

// Activate triggers the activation of a dynamic element
func (s *RawDynamicElement) Activate() error {
	s.active = true
	return nil
}

// Reset resets the state of an dynamic element
func (s *RawDynamicElement) Reset() error {
	s.active = false
	return nil
}

// IsActive returns the activation status of a dynamic element
func (s *RawDynamicElement) IsActive() bool {
	return s.active
}

// NewDynamicElement creates a new dynamic element
func NewDynamicElement(rect *f32.Rectangle) *RawDynamicElement {
	return &RawDynamicElement{
		Area: gui.Rect(rect),
		halo: 4,
	}
}

func checkRect(rect f32.Rectangle, p f32.Point, s float32) bool {
	r := f32.Rectangle{
		Min: f32.Point{X: p.X - s, Y: p.Y - s},
		Max: f32.Point{X: p.X + s, Y: p.Y + s},
	}
	return !rect.Intersect(r).Empty()
}
