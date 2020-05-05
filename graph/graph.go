package graph

import (
	"gioui.org/layout"
	oremi "github.com/drakos74/oremi/internal"
	"github.com/drakos74/oremi/internal/data/model"
)

const (
	width  = 1200
	height = 800
)

type Collection interface {
	Title() string
	Add(labels []string, x ...float64)
	Series() model.Collection
}

type RawCollection struct {
	title  string
	labels []string
	series *model.Series
}

func New(title string, labels ...string) *RawCollection {
	return &RawCollection{
		title:  title,
		series: model.NewSeries(labels...),
	}
}

func (r *RawCollection) Add(labels []string, x ...float64) {
	r.series.Add(model.NewVector(labels, x...))
}

func (r RawCollection) Title() string {
	return r.title
}

func (r RawCollection) Series() model.Collection {
	return r.series
}

func Draw(title string, collections ...Collection) {

	ccs := make(map[string]map[string]model.Collection)

	for _, collection := range collections {
		ccs[collection.Title()] = map[string]model.Collection{collection.Title(): collection.Series()}
	}

	oremi.Draw(title, layout.Vertical, width, height, ccs)
}
