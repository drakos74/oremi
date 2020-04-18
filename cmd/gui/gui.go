package main

import (
	"github/drakos74/oremi/internal/gui"
	"github/drakos74/oremi/internal/gui/canvas/entity"
)

func main() {
	var scene gui.Scene
	scene.WithDimensions(800, 600)

	scene.Add(entity.NewRect(50, 50, 600, 400))

	scene.Run()
}
