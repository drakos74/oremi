package oremi

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	datamodel "github.com/drakos74/oremi/internal/data/model"
	"github.com/drakos74/oremi/internal/gui/canvas"
	uimodel "github.com/drakos74/oremi/internal/gui/model"
	"github.com/drakos74/oremi/internal/gui/style"

	"gioui.org/layout"

	"github.com/drakos74/oremi/internal/gui"
	entity "github.com/drakos74/oremi/internal/gui/canvas/graph"

	"gioui.org/f32"
)

func Draw(title string, axis layout.Axis, width, height float32, collection map[string]map[string]datamodel.Collection) {

	cs := len(collection)

	scene := gui.New().
		WithTitle(title).
		WithDimensions(width+(float32(cs)*(gui.Inset+10)), height+(float32(cs)*(gui.Inset+10)))

	// TODO : fix the layout and collection widths/heights properly
	w := 2*width - 600
	h := 2*height - 50

	switch axis {
	case layout.Horizontal:
		w = w / float32(cs)
	case layout.Vertical:
		h = h / float32(cs)
	}

	graphView := gui.NewView(axis)

	i := 0
	controllers := make([]canvas.Control, 0)
	for title, cc := range collection {
		g := f32.Rectangle{
			Min: f32.Point{X: gui.Inset, Y: gui.Inset},
			Max: f32.Point{X: w, Y: h},
		}

		c, l := filterCollections(cc, datamodel.Size)

		if len(c) > 0 {
			graph := entity.NewChart(l, &g)
			for subtitle, c := range c {
				// TODO : unify building of controls with the collection call
				controller := graph.AddCollection(fmt.Sprintf("%s-%s", title, subtitle), uimodel.NewSeries(c), true)
				controllers = append(controllers, controller)
			}
			graphView.Add(graph)
			i++
		}
	}

	// TODO : layout the autosuggest and checkbox group at the top
	autoSuggestInput := style.NewInput()

	go func() {

		exec := func(cnl chan struct{}, e func()) {
			select {
			case <-time.NewTicker(66 * time.Millisecond).C:
				e()
			case <-cnl:
				return
			}
		}

		cnl := make(chan struct{})

		for event := range autoSuggestInput.Trigger() {
			select {
			case cnl <- struct{}{}:
			default:
			}
			// TODO: fix the acknowledgement path
			//controller.Ack() <- canvas.Ack
			go exec(cnl, func() {
				text := strings.ReplaceAll(event.S, " ", "(.*?)")
				for _, controller := range controllers {
					// TODO : create a proper separate interface for this actions
					if match, _ := regexp.MatchString(text, controller.Label()); !match {
						controller.Disable()
					} else {
						controller.Enable()
					}
				}
			})
		}

	}()

	autoSuggest := gui.NewView(layout.Vertical).WithMaxHeight(30)
	autoSuggest.Add(autoSuggestInput)

	controllerView := gui.NewView(layout.Vertical)
	controllerView.Add(autoSuggest)
	controlView := gui.NewView(layout.Vertical).WithMaxHeight(height + gui.Inset)
	controllerView.Add(controlView)

	screenView := gui.NewView(layout.Horizontal).WithMaxHeight(height + gui.Inset)

	group := style.NewCheckboxControlGroup(true, controllers...)
	controlView.Add(group)
	for _, controller := range controllers {
		controlView.Add(controller)

	}

	screenView.Add(graphView, controllerView)

	scene.Add(screenView)

	scene.Run()
}

func filterCollections(collections map[string]datamodel.Collection, filter datamodel.Filter) (map[string]datamodel.Collection, []string) {
	cc := make(map[string]datamodel.Collection)
	var labels []string
	for key, collection := range collections {
		if filter(collection) {
			cc[key] = collection
			// TODO : be more strict on the labels
			labels = collection.Labels()
		}
	}
	return cc, labels
}
