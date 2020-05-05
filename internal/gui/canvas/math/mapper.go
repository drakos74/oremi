package math

import (
	"fmt"
	"math"

	"gioui.org/f32"
)

type Transform func(x float32) float32

type Mapper interface {
	ScaleX() Transform
	DeScaleX() Transform
	ScaleY() Transform
	DeScaleY() Transform
}

type StackedMapper struct {
	stack []Mapper
}

func NewStackedMapper(stack ...Mapper) *StackedMapper {
	return &StackedMapper{stack: stack}
}

func (s StackedMapper) ScaleX() Transform {
	return func(x float32) float32 {
		for _, s := range s.stack {
			x = s.ScaleX()(x)
		}
		return x
	}
}

func (s StackedMapper) DeScaleX() Transform {
	return func(sx float32) float32 {
		for _, s := range s.stack {
			sx = s.DeScaleX()(sx)
		}
		return sx
	}
}

func (s StackedMapper) ScaleY() Transform {
	return func(y float32) float32 {
		for _, s := range s.stack {
			y = s.ScaleY()(y)
		}
		return y
	}
}

func (s StackedMapper) DeScaleY() Transform {
	return func(sy float32) float32 {
		for _, s := range s.stack {
			sy = s.DeScaleY()(sy)
		}
		return sy
	}
}

type VoidCalcMapper struct {
}

func (v VoidCalcMapper) ScaleX() Transform {
	return func(x float32) float32 {
		return x
	}
}

func (v VoidCalcMapper) DeScaleX() Transform {
	return func(x float32) float32 {
		return x
	}
}

func (v VoidCalcMapper) ScaleY() Transform {
	return func(x float32) float32 {
		return x
	}
}

func (v VoidCalcMapper) DeScaleY() Transform {
	return func(x float32) float32 {
		return x
	}
}

type CoordinateMapper struct {
	scale float32
	rect  *f32.Rectangle
}

func (c CoordinateMapper) ScaleX() Transform {
	return func(x float32) float32 {
		return scaleX(*c.rect, c.scale, x)
	}
}

func (c CoordinateMapper) DeScaleX() Transform {
	return func(sx float32) float32 {
		return deScaleX(*c.rect, c.scale, sx)
	}
}

func (c CoordinateMapper) ScaleY() Transform {
	return func(y float32) float32 {
		return scaleY(*c.rect, c.scale, y)
	}
}

func (c CoordinateMapper) DeScaleY() Transform {
	return func(sy float32) float32 {
		return deScaleY(*c.rect, c.scale, sy)
	}
}

// LinearMapper scales values for x and y linearly to certain ranges and vice versa
type LinearMapper struct {
	min   *f32.Point
	max   *f32.Point
	scale float32
}

func (l LinearMapper) String() string {
	return fmt.Sprintf("%v", struct {
		min   f32.Point
		max   f32.Point
		scale float32
	}{
		min:   *l.min,
		max:   *l.max,
		scale: l.scale,
	})
}

func (l *LinearMapper) Max(pmax f32.Point) bool {
	var recalc bool
	if pmax.X > l.max.X {
		l.max.X = pmax.X
		recalc = true
	}
	if pmax.Y > l.max.Y {
		l.max.Y = pmax.Y
		recalc = true
	}
	return recalc
}

func (l *LinearMapper) Min(pmin f32.Point) bool {
	var recalc bool
	if pmin.X < l.min.X {
		l.min.X = pmin.X
		recalc = true
	}
	if pmin.Y < l.min.Y {
		l.min.Y = pmin.Y
		recalc = true
	}
	return recalc
}

// NewLinearMapper creates a new linearly scale calculation element
func NewLinearMapper(scale float32) *LinearMapper {
	return newLinearMapper(scale, f32.Point{
		X: math.MaxFloat32,
		Y: math.MaxFloat32,
	}, f32.Point{
		X: 0,
		Y: 0,
	})
}

// NewLinearMapper creates a new linearly scale calculation element
func newLinearMapper(scale float32, min, max f32.Point) *LinearMapper {
	return &LinearMapper{
		min:   &min,
		max:   &max,
		scale: scale,
	}
}

func (l LinearMapper) DeScaleX() Transform {
	return func(x float32) float32 {
		return x/l.scale*(l.max.X-l.min.X) + l.min.X
	}
}

func (l LinearMapper) ScaleX() Transform {
	return func(sx float32) float32 {
		return l.scale * (sx - l.min.X) / safe(l.max.X-l.min.X)
	}
}

func (l LinearMapper) DeScaleY() Transform {
	return func(y float32) float32 {
		return y/l.scale*(l.max.Y-l.min.Y) + l.min.Y
	}
}

func (l LinearMapper) ScaleY() Transform {
	return func(sy float32) float32 {
		return l.scale * (sy - l.min.Y) / safe(l.max.Y-l.min.Y)
	}
}

func NewRawCalcElement(rect *f32.Rectangle, scale float32) *CoordinateMapper {
	return &CoordinateMapper{
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
	return (value - rect.Min.X) / safe(rect.Max.X-rect.Min.X) * scale
}

// scaleY calculates the 'real' y - coordinate of a relative value to the grid
func scaleY(rect f32.Rectangle, scale, value float32) float32 {
	return rect.Max.Y - ((rect.Max.Y - rect.Min.Y) * value / scale)
}

// deScaleY calculates the 'relative' y - coordinate of a 'real' value
func deScaleY(rect f32.Rectangle, scale, value float32) float32 {
	return scale - (value-rect.Min.Y)/safe(rect.Max.Y-rect.Min.Y)*scale
}

// safe makes sure that we dont encounter NaN when dividing by '0'
func safe(f float32) float32 {
	if f == 0 {
		return 1
	}
	return f
}

// Float32 converts to a float32 and panics if there is loss of precision
func Float32(f float64) float32 {
	x := float32(f)
	if float64(x) != f {
		panic(fmt.Sprintf("precision loss for %f vs %f", f, x))
	}
	return x
}
