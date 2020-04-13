package entity

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"github.com/google/uuid"
	"oremi/internal/canvas"
	"oremi/internal/model"
)

const scale = 1000

type Graph struct {
	canvas.RawElement
	Scene
	rect        *f32.Rectangle
	max         *f32.Point
	collections map[uint32]model.Collection
	points      map[uint32][]uint32
}

func NewGraph(rect *f32.Rectangle) *Graph {
	g := &Graph{
		*canvas.NewRawElement(),
		*NewScene(rect),
		rect,
		&f32.Point{
			X: 0,
			Y: 0,
		},
		make(map[uint32]model.Collection),
		make(map[uint32][]uint32),
	}
	g.AxisX()
	g.AxisY()
	return g
}

func (g *Graph) Event(e *pointer.Event) (bool, error) {
	if g.Scene.IsActive() {
		p := f32.Point{
			X: g.deScaleX(e.Position.X),
			Y: g.deScaleY(e.Position.Y),
		}
		println(fmt.Sprintf("cursor=%v", p))
	}

	return g.Scene.Event(e)
}

func (g *Graph) Point(label string, p f32.Point) uint32 {
	sp := f32.Point{
		X: g.scaleX(p.X),
		Y: g.scaleY(p.Y),
	}
	point := NewPoint(label, sp)
	g.Scene.RawCompoundElement.Add(point)
	return point.ID()
}

func (g *Graph) AxisX() {
	so := f32.Point{
		X: g.scaleX(0),
		Y: g.scaleY(0),
	}
	g.Scene.RawCompoundElement.Add(NewAxisX(so, g.rect.Max.X-g.rect.Min.X, 10))
}

func (g *Graph) AxisY() {
	so := f32.Point{
		X: g.scaleX(0),
		Y: g.scaleY(scale),
	}
	g.Scene.RawCompoundElement.Add(NewAxisY(so, g.rect.Max.Y-g.rect.Min.Y, 10))
}

// scaleX calculates the 'real' x - coordinate of a relative value to the grid
func (g *Graph) scaleX(value float32) float32 {
	return g.rect.Min.X + ((g.rect.Max.X - g.rect.Min.X) * value / scale)
}

// deScaleX calculates the 'relative' x - coordinate of a 'real' value
func (g *Graph) deScaleX(value float32) float32 {
	return (value - g.rect.Min.X) / (g.rect.Max.X - g.rect.Min.X) * scale
}

// scaleY calculates the 'real' y - coordinate of a relative value to the grid
func (g *Graph) scaleY(value float32) float32 {
	return g.rect.Max.Y - ((g.rect.Max.Y - g.rect.Min.Y) * value / scale)
}

// deScaleX calculates the 'relative' y - coordinate of a 'real' value
func (g *Graph) deScaleY(value float32) float32 {
	return scale - (value-g.rect.Min.Y)/(g.rect.Max.Y-g.rect.Min.Y)*scale
}

// computation specific methods

func (g *Graph) AddCollection(series model.Collection) {
	// TODO : we assume here that minimum is always '0'.
	// BUT : we should handle also negative values
	_, max := series.Edge()
	var recalcNeeded bool
	if float32(max.Coords()[0]) > g.max.X {
		g.max.X = float32(max.Coords()[0])
		recalcNeeded = true
	}
	if float32(max.Coords()[1]) > g.max.Y {
		g.max.Y = float32(max.Coords()[1])
		recalcNeeded = true
	}
	if recalcNeeded {
		for sId, c := range g.collections {
			delete(g.collections, sId)
			for _, pId := range g.points[sId] {
				g.Remove(pId)
			}
			delete(g.points, sId)
			g.eval(sId, c)
		}
	}
	g.eval(uuid.New().ID(), series)
}

func (g *Graph) eval(seriesId uint32, series model.Collection) {

	var points = make([]uint32, series.Size())
	i := 0
	for {
		element, hasNext := series.Next()
		if element != nil {
			x := float32(element.Coords()[0]) / g.max.X
			y := float32(element.Coords()[1]) / g.max.Y
			id := g.Point(
				fmt.Sprintf("%f - %f", element.Coords()[0], element.Coords()[1]),
				f32.Point{
					X: scale * x,
					Y: scale * y,
				})
			points[i] = id
		}
		if !hasNext {
			break
		}
		i++
	}
	g.collections[seriesId] = series
	g.points[seriesId] = points
}
