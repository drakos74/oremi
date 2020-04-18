package bench

import (
	"github/drakos74/oremi/internal/data/model"
	"github/drakos74/oremi/internal/gui"
	"github/drakos74/oremi/internal/gui/canvas/entity"
	uimodel "github/drakos74/oremi/internal/gui/model"

	"gioui.org/f32"
)

// DrawCollection draws a collection of data points
func DrawCollection(collection model.Collection, width, height float32) {
	var scene gui.Scene
	scene.WithDimensions(width, height)

	g := f32.Rectangle{
		Min: f32.Point{X: 50, Y: 50},
		Max: f32.Point{X: width * 2 * 95 / 100, Y: height * 2 * 95 / 100},
	}
	graph := entity.NewGraph(&g)
	graph.AddCollection(uimodel.NewSeries(collection))

	scene.Add(graph)
	scene.Run()
}
