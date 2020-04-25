package main

import (
	oremi "github/drakos74/oremi/internal"
	"github/drakos74/oremi/internal/data/model"
	"github/drakos74/oremi/internal/data/source/generator"
)

func main() {

	oremi.DrawScene("math", 1200, 800,
		map[string][]model.Collection{
			"mathematical-functions": {
				generator.NewPolynomial(120, 0, 0.1, 0, 1),
				generator.NewLine(200, 2, 0, 0.1),
				generator.NewLine(200, 1, 0, 0.1),
				generator.NewExponential(500, 1, 0.01),
			},
		},
	)

}
