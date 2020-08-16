package graph

import (
	"fmt"

	"gioui.org/layout"
	"github.com/drakos74/oremi"
	"github.com/drakos74/oremi/internal/data/model"
)

const (
	width  = 1200
	height = 800
)

type Collection interface {
	Title() string
	Add(index string, labels []string, x ...float64)
	NewSeries(index string, labels ...string)
}

type RawCollection struct {
	title  string
	labels []string
	series map[string]*model.Series
}

func New(title string) *RawCollection {
	return &RawCollection{
		title:  title,
		series: make(map[string]*model.Series),
	}
}

func (r *RawCollection) NewSeries(index string, labels ...string) {
	r.series[index] = model.NewSeries(labels...)
}

func (r *RawCollection) Add(index string, labels []string, x ...float64) {
	if _, ok := r.series[index]; !ok {
		panic(fmt.Sprintf("unknown series with index %v", index))
	}
	r.series[index].Add(model.NewVector(labels, x...))
}

func (r RawCollection) Title() string {
	return r.title
}

func (r RawCollection) Series() map[string]model.Collection {
	collections := make(map[string]model.Collection)
	for index, series := range r.series {
		collections[index] = series
	}
	return collections
}

func Draw(title string, collections ...*RawCollection) {

	ccs := make(map[string]map[string]oremi.Collection)

	colors := oremi.Palette(10)
	for _, collection := range collections {
		collections := make(map[string]oremi.Collection)
		for k, col := range collection.Series() {
			collections[k] = oremi.New(col).Color(colors.Get(k))
		}
		ccs[collection.Title()] = collections
	}

	oremi.Draw(title, layout.Vertical, width, height, ccs)
}
