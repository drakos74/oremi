package gui

import (
	"fmt"
	"github/drakos74/oremi/internal/gui/canvas"
	"github/drakos74/oremi/internal/gui/canvas/entity"
	"log"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/widget/material"

	"gioui.org/app"
	"gioui.org/unit"
)

type Scene struct {
	width    float32
	height   float32
	elements []canvas.Element
}

func (s *Scene) WithDimensions(width, height float32) {
	s.width = width
	s.height = height
}

func (s *Scene) Add(element canvas.Element) {
	s.elements = append(s.elements, element)
}

func (s *Scene) Run() {
	go func() {
		w := app.NewWindow(app.Title("graph-demo"), app.Size(unit.Dp(s.width), unit.Dp(s.height)))
		if err := loop(s, w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(screen *Scene, w *app.Window) error {

	th := material.NewTheme()
	gtx := layout.NewContext(w.Queue())

	scene := entity.NewContainer(&f32.Rectangle{
		Min: f32.Point{0, 0},
		// TODO : we need to get rid of this inconsistency at some point e.g. of multiplying by 2 ...
		Max: f32.Point{2 * screen.width, 2 * screen.height},
	})

	for _, element := range screen.elements {
		println(fmt.Sprintf("element = %v", element))
		scene.Add(element)
	}

	for e := range w.Events() {
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

	return nil
}
