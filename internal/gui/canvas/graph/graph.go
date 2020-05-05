package graph

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/drakos74/oremi/internal/gui/canvas"
	uimath "github.com/drakos74/oremi/internal/gui/canvas/math"
	"github.com/drakos74/oremi/internal/gui/model"
	"github.com/drakos74/oremi/internal/gui/style"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/google/uuid"
)

const scale = 1000

// Chart is a graph object designed to hold all the graph contents as child elements
type Chart struct {
	uimath.CoordinateMapper
	canvas.Container
	scale       *uimath.LinearMapper
	collections map[uint32]collection
	controllers map[uint32]canvas.Control
	points      map[uint32][]uint32
	xaxis       axis
	yaxis       axis
	labels      []string
	trigger     canvas.Events
}

type axis struct {
	axis       uint32
	delimiters []uint32
}

type collection struct {
	model.Collection
	title string
}

// NewChart creates a new graph
func NewChart(labels []string, rect *f32.Rectangle) *Chart {

	if len(labels) < 2 {
		log.Fatalf("cannot draw 2-d graph with only one dimension: %v", labels)
	}

	uiCoordinates := uimath.NewRawCalcElement(rect, scale)
	dataCoordinates := uimath.NewLinearMapper(scale)

	g := &Chart{
		*uiCoordinates,
		*canvas.NewContainer(rect),
		dataCoordinates,
		make(map[uint32]collection),
		make(map[uint32]canvas.Control),
		make(map[uint32][]uint32),
		axis{},
		axis{},
		labels,
		make(canvas.Events),
	}

	g.Axis()

	//g.Add(style.NewCheckBox())
	return g
}

func (g *Chart) Axis() {
	// TODO : we should make the labels flexible and connected to the appropriate dimensions of the vectors
	xaxis, xDelim := g.AxisX(g.labels[0])
	yaxis, yDelim := g.AxisY(g.labels[1])

	xaxisDelims := make([]uint32, len(xDelim))
	yaxisDelims := make([]uint32, len(yDelim))

	for i, xd := range xDelim {
		xaxisDelims[i] = xd.ID()
	}

	for j, yd := range yDelim {
		yaxisDelims[j] = yd.ID()
	}
	g.xaxis = axis{
		axis:       xaxis.ID(),
		delimiters: xaxisDelims,
	}
	g.yaxis = axis{
		axis:       yaxis.ID(),
		delimiters: yaxisDelims,
	}
}

func (g *Chart) RemoveAxis() {
	xaxis := g.xaxis
	g.Remove(xaxis.axis)
	for _, xd := range xaxis.delimiters {
		g.Remove(xd)
	}

	yaxis := g.yaxis
	g.Remove(yaxis.axis)
	for _, yd := range yaxis.delimiters {
		g.Remove(yd)
	}
}

func (g *Chart) Draw(gtx *layout.Context, th *material.Theme) error {
	select {
	case <-g.trigger:
		g.Refresh()
	default:
		//nothing to process
	}
	return g.Container.Draw(gtx, th)
}

// Point adds a point to the graph
func (g *Chart) Point(label string, p f32.Point, control canvas.Control) uint32 {
	sp := f32.Point{
		X: g.ScaleX()(p.X),
		Y: g.ScaleY()(p.Y),
	}
	point := NewPoint(label, sp)
	g.Add(point, control)
	return point.ID()
}

// AxisX adds an x axis to the graph
func (g *Chart) AxisX(label string) (*Axis, []*Delimiter) {
	so := f32.Point{
		X: g.ScaleX()(0),
		Y: g.ScaleY()(0),
	}
	// TODO : fix the calcElement parameter to take into account the max
	rect := g.Rect()
	xAxis := NewAxisX(label, so, rect.Max.X-rect.Min.X)
	g.Add(xAxis, nil)
	delimiters := xAxis.Delimiters(10, uimath.NewStackedMapper(g.CoordinateMapper, g.scale))
	for _, d := range delimiters {
		g.Add(d, nil)
	}
	return xAxis, delimiters
}

// AxisY adds a y axis to the graph
func (g *Chart) AxisY(label string) (*Axis, []*Delimiter) {
	so := f32.Point{
		X: g.ScaleX()(0),
		Y: g.ScaleY()(scale),
	}
	// TODO : fix the calcElement parameter to take into account the max
	rect := g.Rect()
	yAxis := NewAxisY(label, so, rect.Max.Y-rect.Min.Y)
	g.Add(yAxis, nil)
	delimiters := yAxis.Delimiters(10, uimath.NewStackedMapper(g.CoordinateMapper, g.scale))
	for _, d := range delimiters {
		g.Add(d, nil)
	}
	return yAxis, delimiters
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
func (g *Chart) AddCollection(title string, col model.Collection, active bool) canvas.Control {
	err := g.fitsModel(col)
	if err != nil {
		log.Fatalf("cannot add collection to graph: %v", err)
	}

	bound := col.Bounds()

	newMax := g.scale.Max(bound.Max)
	newMin := g.scale.Min(bound.Min)

	if newMax || newMin {
		// update the existing collections in terms of scaling
		for sId, c := range g.collections {
			// NOTE : this is a tricky point ... need to handle with care
			g.remove(sId)
			g.add(sId, c.title, c, g.controllers[sId])
		}
	}

	controller := style.NewCheckBox(title, active)
	sId := uuid.New().ID()
	g.add(sId, title, col, controller)
	g.collections[sId] = collection{
		col,
		title,
	}
	g.controllers[sId] = controller

	// TODO : abstract the default trigger receiver logic into dedicated interface e.g. Loader with Refresh
	go func() {

		exec := func(cnl chan struct{}, e func()) {
			select {
			case <-time.NewTicker(66 * time.Millisecond).C:
				e()
			case <-cnl:
				return
			}
		}

		cnl := make(chan struct{})

		for event := range controller.Trigger() {
			select {
			case cnl <- struct{}{}:
			default:
			}
			// TODO: fix the acknowledgement path
			//controller.Ack() <- canvas.Ack
			go exec(cnl, func() {
				g.trigger <- event
			})
		}

	}()

	return controller
}

func (g *Chart) Refresh() {
	// TODO : make this dynamic ready e.g. redrawn elements should keep their events
	g.RemoveAxis()

	g.scale = uimath.NewLinearMapper(scale)

	for id, collection := range g.collections {
		if g.controllers[id].IsActive() {
			bound := collection.Bounds()
			g.scale.Max(bound.Max)
			g.scale.Min(bound.Min)
		}
	}

	g.Axis()

	for sId, c := range g.collections {
		// NOTE : this is a tricky point ... need to handle with care
		g.remove(sId)
		g.add(sId, c.title, c, g.controllers[sId])
	}

}

// remove removes a collection and it's points
func (g *Chart) remove(sId uint32) {
	for _, pId := range g.points[sId] {
		g.Remove(pId)
		g.Remove(pId)
	}
	delete(g.points, sId)
}

// TODO : allow for dynamically changing collections e.g. adding data points on the graph
// add scales the model series into canvas coordinates scale
func (g *Chart) add(sId uint32, title string, collection model.Collection, controller canvas.Control) {
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
				}, controller)
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
