package graph

import (
	"fmt"
	"github/drakos74/oremi/internal/gui/canvas"
	uimath "github/drakos74/oremi/internal/gui/canvas/math"
	"github/drakos74/oremi/internal/gui/model"
	"log"
	"math"
	"strings"

	"gioui.org/f32"
	"github.com/google/uuid"
)

const scale = 1000

// Chart is a graph object designed to hold all the graph contents as child elements
type Chart struct {
	uimath.CoordinateMapper
	canvas.Container
	scale       *uimath.LinearMapper
	rect        *f32.Rectangle
	collections map[uint32]model.Collection
	points      map[uint32][]uint32
	labels      []string
}

// NewChart creates a new graph
func NewChart(labels []string, rect *f32.Rectangle) *Chart {

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

	g := &Chart{
		*uiCoordinates,
		*canvas.NewContainer(rect),
		dataCoordinates,
		rect,
		make(map[uint32]model.Collection),
		make(map[uint32][]uint32),
		labels,
	}
	// TODO : we should make the labels flexible and connected to the appropriate dimensions of the vectors
	g.AxisX(labels[0])
	g.AxisY(labels[1])
	//g.AddItem(style.NewCheckBox())
	return g
}

// Point adds a point to the graph
func (g *Chart) Point(label string, p f32.Point) uint32 {
	sp := f32.Point{
		X: g.ScaleX()(p.X),
		Y: g.ScaleY()(p.Y),
	}
	point := NewPoint(label, sp)
	g.Add(point)
	return point.ID()
}

// AxisX adds an x axis to the graph
func (g *Chart) AxisX(label string) {
	so := f32.Point{
		X: g.ScaleX()(0),
		Y: g.ScaleY()(0),
	}
	// TODO : fix the calcElement parameter to take into account the max
	xAxis := NewAxisX(label, so, g.rect.Max.X-g.rect.Min.X)
	g.Add(xAxis)
	delimiters := xAxis.Delimiters(10, uimath.NewStackedMapper(g.CoordinateMapper, g.scale))
	for _, d := range delimiters {
		g.Add(d)
	}
}

// AxisY adds a y axis to the graph
func (g *Chart) AxisY(label string) {
	so := f32.Point{
		X: g.ScaleX()(0),
		Y: g.ScaleY()(scale),
	}
	// TODO : fix the calcElement parameter to take into account the max
	yAxis := NewAxisY(label, so, g.rect.Max.Y-g.rect.Min.Y)
	g.Add(yAxis)
	delimiters := yAxis.Delimiters(10, uimath.NewStackedMapper(g.CoordinateMapper, g.scale))
	for _, d := range delimiters {
		g.Add(d)
	}
}

// model validation methods
func (g *Chart) fitsModel(collection model.Collection) error {
	for i, label := range collection.Labels() {
		if g.labels[i] != label {
			return fmt.Errorf("model inconsistency on labels %v vs %v", g.labels, collection.Labels())
		}
	}
	return nil
}

// computation specific methods

// AddCollection adds a series model collection to the graph
func (g *Chart) AddCollection(collection model.Collection) {
	err := g.fitsModel(collection)
	if err != nil {
		log.Fatalf("cannot add collection to graph: %v", err)
	}

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
func (g *Chart) remove(sId uint32) {
	for _, pId := range g.points[sId] {
		g.Remove(pId)
	}
	delete(g.points, sId)
}

// add scales the model series into canvas coordinates scale
func (g *Chart) add(sId uint32, collection model.Collection) {
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
