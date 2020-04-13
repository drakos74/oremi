package generator

import "github/drakos74/oremi/internal/model"

type Y func(x float64) float64

type Generator struct {
	f []Y
}

func (g *Generator) Next(p model.Element) model.Element {
	coords := make([]float64, p.Dim())
	for i, c := range p.Coords() {
		coords[i] = g.f[i](c)
	}
	return model.NewPoint("", coords...)
}
