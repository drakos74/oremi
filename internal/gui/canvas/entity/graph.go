package entity

import (
	"fmt"
	"github/drakos74/oremi/internal/gui/canvas"
	"github/drakos74/oremi/internal/gui/model"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"github.com/google/uuid"
)

const scale = 1000

// Graph is a graph object designed to hold all the graph contents as child elements
type Graph struct {
	canvas.RawElement
	Container
	rect        *f32.Rectangle
	max         *f32.Point
	collections map[uint32]model.Collection
	points      map[uint32][]uint32
}

// NewGraph creates a new graph
func NewGraph(rect *f32.Rectangle) *Graph {
	g := &Graph{
		*canvas.NewRawElement(),
		*NewContainer(rect),
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

// Event propagates the events to all child elements of the graph
func (g *Graph) Event(e *pointer.Event) (bool, error) {
	if g.Container.IsActive() {
		p := f32.Point{
			X: g.deScaleX(e.Position.X),
			Y: g.deScaleY(e.Position.Y),
		}
		// TODO : remove or enable only with appropriate config or action
		println(fmt.Sprintf("cursor=%v", p))
	}

	return g.Container.Event(e)
}

// Point adds a point to the graph
func (g *Graph) Point(label string, p f32.Point) uint32 {
	sp := f32.Point{
		X: g.scaleX(p.X),
		Y: g.scaleY(p.Y),
	}
	point := NewPoint(label, sp)
	g.Add(point)
	return point.ID()
}

// AxisX adds an x axis to the graph
func (g *Graph) AxisX() {
	so := f32.Point{
		X: g.scaleX(0),
		Y: g.scaleY(0),
	}
	g.Add(NewAxisX(so, g.rect.Max.X-g.rect.Min.X, 10))
}

// AxisY adds a y axis to the graph
func (g *Graph) AxisY() {
	so := f32.Point{
		X: g.scaleX(0),
		Y: g.scaleY(scale),
	}
	g.Add(NewAxisY(so, g.rect.Max.Y-g.rect.Min.Y, 10))
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

// AddCollection adds a series model collection to the graph
func (g *Graph) AddCollection(collection model.Collection) {
	// TODO : we assume here that minimum is always '0'.
	// BUT : we should handle also negative values
	bound := collection.Bounds()
	var doRecalc bool
	if bound.Max.X > g.max.X {
		g.max.X = bound.Max.X
		doRecalc = true
	}
	if bound.Max.Y > g.max.Y {
		g.max.Y = bound.Max.Y
		doRecalc = true
	}

	if doRecalc {
		for sId, c := range g.collections {
			g.remove(sId)
			g.add(sId, c)
		}
	}

	sId := uuid.New().ID()
	g.add(sId, collection)
	g.collections[sId] = collection

}

// remove removes a collection and it's points
func (g *Graph) remove(sId uint32) {
	for _, pId := range g.points[sId] {
		g.Remove(pId)
	}
	delete(g.points, sId)
}

// add scales the model series into canvas coordinates scale
func (g *Graph) add(sId uint32, collection model.Collection) {

	collection.Reset()
	var points = make([]uint32, collection.Size())
	i := 0
	for {
		point, ok, hasNext := collection.Next()
		if ok {
			x := point.X / g.max.X
			y := point.Y / g.max.Y
			id := g.Point(
				point.Label,
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
	g.points[sId] = points
}
