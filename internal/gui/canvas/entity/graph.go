package entity

import (
	"errors"
	"fmt"
	"github/drakos74/oremi/internal/gui/canvas"
	uimath "github/drakos74/oremi/internal/gui/canvas/math"
	"github/drakos74/oremi/internal/gui/model"
	"log"
	"math"
	"strings"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"github.com/google/uuid"
)

const scale = 1000

// Graph is a graph object designed to hold all the graph contents as child elements
type Graph struct {
	canvas.RawElement
	uimath.CoordinateMapper
	Container
	scale       *uimath.LinearMapper
	rect        *f32.Rectangle
	collections map[uint32]model.Collection
	points      map[uint32][]uint32
	labels      []string
}

// NewGraph creates a new graph
func NewGraph(labels []string, rect *f32.Rectangle) *Graph {

	if len(labels) < 2 {
		log.Fatalf("cannot draw 2-d graph with only one dimension: %v", labels)
	}

	uiCoordinates := uimath.NewRawCalcElement(rect, scale)

	dataCoordinates := uimath.NewLinearMapper(scale, f32.Point{
		X: math.MaxFloat32,
		Y: math.MaxFloat32,
	}, f32.Point{
		X: 0,
		Y: 0,
	})

	g := &Graph{
		*canvas.NewRawElement(),
		*uiCoordinates,
		*NewContainer(rect),
		dataCoordinates,
		rect,
		make(map[uint32]model.Collection),
		make(map[uint32][]uint32),
		labels,
	}
	// TODO : we should make the labels flexible and connected to the appropriate dimensions of the vectors
	g.AxisX(labels[0])
	g.AxisY(labels[1])
	return g
}

// Event propagates the events to all child elements of the graph
func (g *Graph) Event(e *pointer.Event) (bool, error) {
	if g.Container.IsActive() {
		p := f32.Point{
			X: g.DeScaleX()(e.Position.X),
			Y: g.DeScaleY()(e.Position.Y),
		}
		// TODO : remove or enable only with appropriate config or action
		println(fmt.Sprintf("cursor=%v", p))
	}

	return g.Container.Event(e)
}

// Point adds a point to the graph
func (g *Graph) Point(label string, p f32.Point) uint32 {
	sp := f32.Point{
		X: g.ScaleX()(p.X),
		Y: g.ScaleY()(p.Y),
	}
	point := NewPoint(label, sp)
	g.Add(point)
	return point.ID()
}

// AxisX adds an x axis to the graph
func (g *Graph) AxisX(label string) {
	so := f32.Point{
		X: g.ScaleX()(0),
		Y: g.ScaleY()(0),
	}
	// TODO : fix the calcElement parameter to take into account the max
	g.Add(NewAxisX(label, so, g.rect.Max.X-g.rect.Min.X, 10, uimath.NewStackedMapper(g.CoordinateMapper, g.scale)))
}

// AxisY adds a y axis to the graph
func (g *Graph) AxisY(label string) {
	so := f32.Point{
		X: g.ScaleX()(0),
		Y: g.ScaleY()(scale),
	}
	// TODO : fix the calcElement parameter to take into account the max
	g.Add(NewAxisY(label, so, g.rect.Max.Y-g.rect.Min.Y, 10, uimath.NewStackedMapper(g.CoordinateMapper, g.scale)))
}

// model validation methods
func (g *Graph) fitsModel(collection model.Collection) error {
	for i, label := range collection.Labels() {
		if g.labels[i] != label {
			return errors.New(fmt.Sprintf("model inconsistency on labels %v vs %v", g.labels, collection.Labels()))
		}
	}
	return nil
}

// computation specific methods

// AddCollection adds a series model collection to the graph
func (g *Graph) AddCollection(collection model.Collection) {
	err := g.fitsModel(collection)
	if err != nil {
		log.Fatalf("cannot add collection to graph: %v", err)
	}

	// TODO : we assume here that minimum is always '0'.
	// TODO : we should handle also negative values
	bound := collection.Bounds()

	newMax := g.scale.Max(bound.Max)
	newMin := g.scale.Min(bound.Min)

	if newMax || newMin {
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
			id := g.Point(
				label(point.Label),
				f32.Point{
					X: g.scale.ScaleX()(point.X),
					Y: g.scale.ScaleY()(point.Y),
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

func label(labels []string) string {
	return strings.Join(labels, "-")
}
