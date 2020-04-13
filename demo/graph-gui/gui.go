package main

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"log"
	"oremi/internal/canvas/entity"
	"oremi/internal/source/generator"

	"gioui.org/app"
)

const (
	width  = 1200
	height = 800
)

// TODO : notice that we need to duplicate this ... seems some inconsistency with the library (?)
var ww = float32(2 * width)
var hh = float32(2 * height)

func main() {
	go func() {
		w := app.NewWindow(app.Title("demo"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {

	th := material.NewTheme()
	gtx := layout.NewContext(w.Queue())

	scene := entity.NewScene(&f32.Rectangle{
		Min: f32.Point{0, 0},
		Max: f32.Point{ww, hh},
	})

	g := f32.Rectangle{
		Min: f32.Point{X: 50, Y: 50},
		Max: f32.Point{X: ww * 95 / 100, Y: hh * 95 / 100},
	}
	graph := entity.NewGraph(&g)

	scene.Add(graph)

	l1 := generator.NewLine()
	l1.Generate(10, 4, 4)
	graph.AddCollection(l1)

	l2 := generator.NewLine()
	l2.Generate(10, 2, 4)
	graph.AddCollection(l2)

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx.Reset(e.Config, e.Size)
				err := scene.Draw(gtx, th)
				if err != nil {
					log.Fatalf("error during scene drawing: %v", err)
				}
				e.Frame(gtx.Ops)
			case pointer.Event:
				active, err := scene.Event(&e)
				if err != nil {
					log.Fatalf("error during event propagation: %v", err)
				}
				// TODO : does not work very well in de-activation conditions
				if active {
					// trigger a new frame event
					w.Invalidate()
				}
			}
		}
	}
}
