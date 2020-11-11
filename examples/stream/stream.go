package main

import (
	"fmt"
	"time"

	"gioui.org/layout"

	"github.com/drakos74/oremi"

	"github.com/drakos74/oremi/internal/data/model"
	"github.com/drakos74/oremi/label"
)

func main() {

	series := model.NewSeries(label.Num("x"), label.Num("y"))

	collection := *oremi.New(series)
	series.Add(model.NewVector([]string{fmt.Sprintf("%d", 0), fmt.Sprintf("%d", 0)}, float64(0), float64(0)))

	go func() {
		for i := 0; i < 100; i++ {
			series.Add(model.NewVector([]string{fmt.Sprintf("%d", i), fmt.Sprintf("%d", i*10)}, float64(i), float64(10*i)))
			time.Sleep(3 * time.Second)
		}
	}()

	oremi.Draw("test", layout.Vertical, 1200, 800, map[string]map[string]oremi.Collection{
		"": {
			"": collection,
		},
	}, 1)

}
