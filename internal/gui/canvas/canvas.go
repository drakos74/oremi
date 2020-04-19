package canvas

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/google/uuid"
)

type Transform func(x float32) float32

// Element is the main abstraction for any object living within the canvas
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

func (s *RawElement) Scale() (width, height float32) {
	return width, height
}

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

// DrawElement represents an elements that can be drawn on the canvas
type DrawElement interface {
	Draw(gtx *layout.Context, th *material.Theme) error
}

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
		if el, ok := element.(DrawElement); ok {
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
		if el, ok := element.(DynamicElement); ok {
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

// Add adds a new element to the group
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

// DynamicElement represents an interactive UI element
type DynamicElement interface {
	Event(e *pointer.Event) (bool, error)
	Activate() error
	Reset() error
	IsActive() bool
}

// RawDynamicElement is the base implementation of a dynamic element
type RawDynamicElement struct {
	active bool
	halo   float32
	rect   f32.Rectangle
}

// Event propagates the scene event to the element
func (s *RawDynamicElement) Event(e *pointer.Event) (bool, error) {
	if !checkRect(s.rect, e.Position, s.halo) {
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
func NewDynamicElement(rect f32.Rectangle) *RawDynamicElement {
	return &RawDynamicElement{
		halo: 4,
		rect: rect,
	}
}

func checkRect(rect f32.Rectangle, p f32.Point, s float32) bool {
	r := f32.Rectangle{
		Min: f32.Point{X: p.X - s, Y: p.Y - s},
		Max: f32.Point{X: p.X + s, Y: p.Y + s},
	}
	return !rect.Intersect(r).Empty()
}

type CalcElement interface {
	ScaleX() Transform
	DeScaleX() Transform
	ScaleY() Transform
	DeScaleY() Transform
}

type RawCalcElement struct {
	scale float32
	rect  *f32.Rectangle
}

func (c RawCalcElement) ScaleX() Transform {
	return func(x float32) float32 {
		return scaleX(*c.rect, c.scale, x)
	}
}

func (c RawCalcElement) DeScaleX() Transform {
	return func(sx float32) float32 {
		return deScaleX(*c.rect, c.scale, sx)
	}
}

func (c RawCalcElement) ScaleY() Transform {
	return func(y float32) float32 {
		return scaleY(*c.rect, c.scale, y)
	}
}

func (c RawCalcElement) DeScaleY() Transform {
	return func(sy float32) float32 {
		return deScaleY(*c.rect, c.scale, sy)
	}
}

func NewRawCalcElement(rect *f32.Rectangle, scale float32) *RawCalcElement {
	return &RawCalcElement{
		scale: scale,
		rect:  rect,
	}
}

// scaleX calculates the 'real' x - coordinate of a relative value to the grid
func scaleX(rect f32.Rectangle, scale, value float32) float32 {
	return rect.Min.X + ((rect.Max.X - rect.Min.X) * value / scale)
}

// deScaleX calculates the 'relative' x - coordinate of a 'real' value
func deScaleX(rect f32.Rectangle, scale, value float32) float32 {
	return (value - rect.Min.X) / (rect.Max.X - rect.Min.X) * scale
}

// scaleY calculates the 'real' y - coordinate of a relative value to the grid
func scaleY(rect f32.Rectangle, scale, value float32) float32 {
	return rect.Max.Y - ((rect.Max.Y - rect.Min.Y) * value / scale)
}

// deScaleY calculates the 'relative' y - coordinate of a 'real' value
func deScaleY(rect f32.Rectangle, scale, value float32) float32 {
	return scale - (value-rect.Min.Y)/(rect.Max.Y-rect.Min.Y)*scale
}
