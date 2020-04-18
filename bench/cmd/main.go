package main

import (
	"flag"
	"fmt"
	"github/drakos74/oremi/bench"
	"github/drakos74/oremi/internal/gui"
	"github/drakos74/oremi/internal/gui/canvas/entity"
	"github/drakos74/oremi/internal/gui/model"
	"log"

	"gioui.org/f32"
)

const (
	width  float32 = 1200
	height float32 = 800
)

func main() {

	file := flag.String("file", "", "benchmark output file")

	flag.Parse()

	println(fmt.Sprintf("parsing benchmark file = %v", *file))

	b, err := bench.ParseBenchmarks(*file)

	if err != nil {
		log.Fatalf("could not parse benchamrks from file '%s': %v", *file, err)
	}

	collection := b.ToCollection()

	var scene gui.Scene
	scene.WithDimensions(width, height)

	g := f32.Rectangle{
		Min: f32.Point{X: 50, Y: 50},
		Max: f32.Point{X: width * 2 * 95 / 100, Y: height * 2 * 95 / 100},
	}
	graph := entity.NewGraph(&g)
	graph.AddCollection(model.NewSeries(collection))

	scene.Add(graph)
	scene.Run()
}
