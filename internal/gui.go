package oremi

import (
	datamodel "github/drakos74/oremi/internal/data/model"
	uimodel "github/drakos74/oremi/internal/gui/model"

	"github/drakos74/oremi/internal/gui"
	"github/drakos74/oremi/internal/gui/canvas/entity"

	"gioui.org/f32"
)

func DrawScene(title string, width, height float32, collection map[string][]datamodel.Collection) {
	var scene gui.Scene
	scene.WithTitle(title)
	scene.WithDimensions(width, height)

	w := width * 2 * 95 / float32(len(collection)*100)

	i := 0
	for _, cc := range collection {
		g := f32.Rectangle{
			Min: f32.Point{X: (float32(i) * w) + 100, Y: 50},
			Max: f32.Point{X: (float32(i) * w) + w, Y: 2 * height * 95 / 100},
		}
		coll := clearCollections(cc)
		if len(coll) > 0 {
			graph := entity.NewGraph(coll[0].Labels(), &g)
			for _, c := range coll {
				graph.AddCollection(uimodel.NewSeries(c))
			}
			scene.Add(graph)
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
