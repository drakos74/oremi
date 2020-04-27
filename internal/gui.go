package oremi

import (
	datamodel "github/drakos74/oremi/internal/data/model"
	uimodel "github/drakos74/oremi/internal/gui/model"

	"gioui.org/layout"

	"github/drakos74/oremi/internal/gui"
	entity "github/drakos74/oremi/internal/gui/canvas/graph"

	"gioui.org/f32"
)

func DrawScene(title string, axis layout.Axis, width, height float32, collection map[string][]datamodel.Collection) {

	cs := len(collection)

	scene := gui.New().
		WithTitle(title).
		WithDimensions(width+(float32(cs)*gui.Inset), height+(float32(cs)*gui.Inset))
	scene.WithLayout(axis)

	// TODO : fix the layout and collection widths/heights properly
	w := 2 * width
	h := 2 * height

	switch axis {
	case layout.Horizontal:
		w = w / float32(cs)
	case layout.Vertical:
		h = h / float32(cs)
	}

	i := 0
	// TODO : fix multiple scene elements (check bench example)
	// TODO : get event / rect border from draw context
	for _, cc := range collection {
		g := f32.Rectangle{
			Min: f32.Point{X: gui.Inset, Y: gui.Inset},
			Max: f32.Point{X: w, Y: h},
		}
		coll := clearCollections(cc)
		if len(coll) > 0 {
			graph := entity.NewChart(coll[0].Labels(), &g)
			for _, c := range coll {
				graph.AddCollection(uimodel.NewSeries(c))
			}
			scene.AddItem(graph)
			i++
		}
	}

	scene.Run()
}

func clearCollections(collections []datamodel.Collection) []datamodel.Collection {
	c := make([]datamodel.Collection, 0)
	for _, col := range collections {
		if col.Size() > 0 {
			c = append(c, col)
		}
	}
	return c
}
