package bench

import (
	"github/drakos74/oremi/internal/data/model"
	"github/drakos74/oremi/internal/gui"
	"github/drakos74/oremi/internal/gui/canvas/entity"
	uimodel "github/drakos74/oremi/internal/gui/model"

	"gioui.org/f32"
)

// DrawCollection draws a collection of data points
func DrawCollections(width, height float32, collection ...model.Collection) {
	var scene gui.Scene
	scene.WithDimensions(width, height)

	cols := clearCollections(collection)

	w := width * 2 * 95 / float32(len(cols)*100)

	for i, c := range cols {
		g := f32.Rectangle{
			Min: f32.Point{X: (float32(i) * w) + 100, Y: 50},
			Max: f32.Point{X: (float32(i) * w) + w, Y: 2 * height * 95 / 100},
		}
		graph := entity.NewGraph(c.Labels()[0], c.Labels()[1], &g)
		graph.AddCollection(uimodel.NewSeries(c))

		scene.Add(graph)
	}

	scene.Run()
}

func clearCollections(collections []model.Collection) []model.Collection {
	c := make([]model.Collection, 0)
	for _, col := range collections {
		if col.Size() > 0 {
			c = append(c, col)
		}
	}
	return c
}
