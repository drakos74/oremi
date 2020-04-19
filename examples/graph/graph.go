package main

import (
	"github/drakos74/oremi/internal/data/source/generator"
	"github/drakos74/oremi/internal/gui"
	"github/drakos74/oremi/internal/gui/canvas/entity"
	"github/drakos74/oremi/internal/gui/model"

	"gioui.org/f32"
)

const (
	width  = 1200
	height = 800
)

// TODO : notice that we need to duplicate this ... seems some inconsistency with the implementation (?)
var ww = float32(2 * width)
var hh = float32(2 * height)

func main() {

	var scene gui.Scene

	scene.WithDimensions(1200, 800)

	g := f32.Rectangle{
		Min: f32.Point{X: 50, Y: 50},
		Max: f32.Point{X: ww * 95 / 100, Y: hh * 95 / 100},
	}
	graph := entity.NewGraph("x", "f(x)", &g)
	graph.AddCollection(model.NewSeries(generator.NewPolynomial(120, 0, 0.1, 0, 1)))
	graph.AddCollection(model.NewSeries(generator.NewLine(200, 2, 0, 0.1)))
	graph.AddCollection(model.NewSeries(generator.NewLine(200, 1, 0, 0.1)))
	graph.AddCollection(model.NewSeries(generator.NewExponential(500, 1, 0.01)))
	scene.Add(graph)

	scene.Run()

}
