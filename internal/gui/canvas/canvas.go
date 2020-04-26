package canvas

import (
	"github/drakos74/oremi/internal/gui"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/google/uuid"
)

// Element is the main abstraction for any object living within the canvas
// TODO : replace usage with item abstraction
type Element interface {
	ID() uint32
	Scale() (width, height float32)
	Offset() (x, y float32)
}

// RawElement is the base implementation for an Element
type RawElement struct {
	id     uint32
	width  float32
	height float32
}

// TODO : make use of Scale or remove
func (s *RawElement) Scale() (width, height float32) {
	return width, height
}

// TODO : make use of Offset or remove
func (s *RawElement) Offset() (x, y float32) {
	return x, y
}

// ID returns the id of the raw element
func (s *RawElement) ID() uint32 {
	return s.id
}

// NewRawElement creates a new raw element
func NewRawElement() *RawElement {
	return &RawElement{
		id:     uuid.New().ID(),
		width:  1000,
		height: 1000,
	}
}

// TODO : replace with items abstraction
// CompoundElement represents an element that can have children
type CompoundElement interface {
	Add(element Element)
	Remove(id uint32)
	Elements(apply Action) (bool, error)
}

// Action defines an action to be applied to an Element
type Action func(element Element) (bool, error)

// DrawFunction is a helper method to invoke the Draw method on an elements
var DrawAction = func(gtx *layout.Context, th *material.Theme) func(element Element) (bool, error) {
	return func(element Element) (bool, error) {
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
var EventAction = func(e *pointer.Event) func(element Element) (bool, error) {
	return func(element Element) (bool, error) {
		var active bool
		// TODO : avoid reflection by keeping the event actions in a slice
		if el, ok := element.(gui.InteractiveItem); ok {
			a, err := el.Event(e)
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
	elements map[uint32]Element
}

// AddItem adds a new element to the group
func (s *RawCompoundElement) Add(element Element) {
	//t := reflect.TypeOf(element)
	// TODO : use to make better use of generic actions, without doing the casting
	s.elements[element.ID()] = element
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
	return &RawCompoundElement{elements: make(map[uint32]Element)}
}
