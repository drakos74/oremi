package gui

import (
	"log"

	"gioui.org/font/gofont"

	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/widget/material"

	"gioui.org/app"
	"gioui.org/unit"
)

func init() {
	gofont.Register()
}

const Inset = 50

// Scene is the main container for the canvas objects
type Scene struct {
	list
	title      string
	width      float32
	height     float32
	layoutList *layout.List
}

func New() *Scene {
	return &Scene{
		list: list{
			items: make(map[uint32]item),
		},
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

func (s *Scene) WithLayout(orientation layout.Axis) {
	s.layoutList = &layout.List{
		Axis: orientation,
	}
}

// live methods

// Render renders all elements in the scene
func (s *Scene) Draw(gtx *layout.Context, th *material.Theme) error {
	// TODO : do recover
	s.layoutList.Layout(gtx, len(s.items), func(i int) {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, s.get(i).draw(gtx, th))
	})
	return nil
}

// Event propagates a scene event
func (s *Scene) Event(e *pointer.Event) (redraw bool, err error) {
	for i := 0; i < len(s.items); i++ {
		if s.get(i).event(e) {
			redraw = true
		}
	}
	return redraw, nil
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

	th := material.NewTheme()
	gtx := layout.NewContext(w.Queue())

	for e := range w.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)
			// TODO : avoid re-drawing if nothing changed
			err := scene.Draw(gtx, th)
			if err != nil {
				// TODO : handle error so that it freezes instead of failing
				// consider using custom error type
				log.Fatalf("could not draw scene: %v", err)
			}
			e.Frame(gtx.Ops)
		case pointer.Event:
			redraw, err := scene.Event(&e)
			if err != nil {
				// TODO : handle error so that it ignores instead of failing
				// consider using custom error type
				log.Fatalf("could not draw scene: %v", err)
			}
			if redraw {
				w.Invalidate()
			}
		}
	}

	return nil
}
