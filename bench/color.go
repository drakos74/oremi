package bench

import (
	"image/color"
	"math/rand"
)

type Color struct {
	colors map[string]color.RGBA
}

func Palette() *Color {
	return &Color{colors: make(map[string]color.RGBA)}
}

func (c *Color) Get(label string) color.RGBA {
	if col, ok := c.colors[label]; ok {
		return col
	}
	cc := color.RGBA{uint8(rand.Intn(125) + 100), uint8(rand.Intn(125) + 100), uint8(rand.Intn(125) + 100), 255}
	c.colors[label] = cc
	return cc
}
