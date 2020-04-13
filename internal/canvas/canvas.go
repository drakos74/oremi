package canvas

import (
	"fmt"
	"reflect"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/google/uuid"
)

type Element interface {
	ID() uint32
}

type RawElement struct {
	id uint32
}

func (s *RawElement) ID() uint32 {
	return s.id
}

func NewRawElement() *RawElement {
	return &RawElement{
		id: uuid.New().ID(),
	}
}

type DrawElement interface {
	Draw(gtx *layout.Context, th *material.Theme) error
}

type CompoundElement interface {
	Add(element Element)
	Remove(id uint32)
	Elements(apply Action) (bool, error)
}

type Action func(element Element) (bool, error)

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

type RawCompoundElement struct {
	elements map[uint32]Element
}

func (s *RawCompoundElement) Add(element Element) {
	t := reflect.TypeOf(element)
	println(fmt.Sprintf("type = %v", t))
	// TODO : use to make better use of generic actions, without doing the casting
	s.elements[element.ID()] = element
}

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

func (s *RawCompoundElement) Remove(id uint32) {
	delete(s.elements, id)
}

func (s *RawCompoundElement) Size() int {
	return len(s.elements)
}

func NewCompoundElement() *RawCompoundElement {
	return &RawCompoundElement{elements: make(map[uint32]Element)}
}

type DynamicElement interface {
	Event(e *pointer.Event) (bool, error)
	Activate() error
	Reset() error
	IsActive() bool
}

type RawDynamicElement struct {
	active bool
	halo   float32
	rect   f32.Rectangle
}

func (s *RawDynamicElement) Event(e *pointer.Event) (bool, error) {
	if !checkRect(s.rect, e.Position, s.halo) {
		err := s.Reset()
		return false, err
	} else {
		err := s.Activate()
		return err == nil, err
	}
}

func (s *RawDynamicElement) Activate() error {
	s.active = true
	return nil
}

func (s *RawDynamicElement) Reset() error {
	s.active = false
	return nil
}

func (s *RawDynamicElement) IsActive() bool {
	return s.active
}

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
