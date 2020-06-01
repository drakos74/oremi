package main

import (
	"github.com/drakos74/oremi"
	"github.com/drakos74/oremi/internal/data/model"
	"github.com/drakos74/oremi/internal/data/source/generator"

	"gioui.org/layout"
)

func main() {

	oremi.Draw("math", layout.Vertical, 1200, 800,
		map[string]map[string]model.Collection{
			"math": {
				"polynomial":  generator.NewPolynomial(120, 0, 0.1, 0, 1),
				"linear-1":    generator.NewLine(200, 2, 0, 0.1),
				"linear-2":    generator.NewLine(200, 1, 0, 0.1),
				"exponential": generator.NewExponential(500, 1, 0.01),
			},
		},
	)

}
