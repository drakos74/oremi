package gui

import (
	"log"

	"gioui.org/io/pointer"

	"gioui.org/font/gofont"
	"gioui.org/widget/material"

	"gioui.org/op"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
)

const Inset = 50

// Scene is the main container for the canvas objects
type Scene struct {
	*View
	title  string
	width  float32
	height float32
}

func New() *Scene {
	return &Scene{
		View: NewView(layout.Vertical),
	}
}

// config methods

// WithTitle defines the window title
func (s *Scene) WithTitle(title string) *Scene {
	s.title = title
	return s
}

// WithDimensions defines the window dimensions
func (s *Scene) WithDimensions(width, height float32) *Scene {
	s.width = width
	s.height = height
	return s
}

// Run start the gui
func (s *Scene) Run() {
	go func() {
		w := app.NewWindow(app.Title(s.title), app.Size(unit.Dp(s.width), unit.Dp(s.height)))
		if err := loop(s, w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(scene *Scene, w *app.Window) error {

	th := material.NewTheme(gofont.Collection())

	var ops op.Ops
	for e := range w.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			// TODO : find a smarter way to 'invalidate'
			w.Invalidate()
			gtx := layout.NewContext(&ops, e)
			// TODO : avoid re-drawing if nothing changed
			_, err := scene.Draw(gtx, th)
			if err != nil {
				// TODO : handle error so that it freezes instead of failing
				// consider using custom error type
				log.Fatalf("could not draw scene: %v", err)
			}
			e.Frame(gtx.Ops)
		case pointer.Event:
			// TODO : dont make use of gtx from this scope
			//redraw, err := scene.Event(&e)
			//if err != nil {
			//	// TODO : handle error so that it ignores instead of failing
			//	// consider using custom error type
			//	log.Fatalf("could not draw scene: %v", err)
			//}
			//if redraw {
			//	w.Invalidate()
			//}
			//default:
			//	println(fmt.Sprintf("e = %+v", e))
		}
	}

	return nil
}
