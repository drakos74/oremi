package bench

import (
	"encoding/binary"
	"hash/fnv"
	"image/color"
)

type Color struct {
	colors map[string]color.RGBA
}

func Palette(num int) *Color {
	colors := Color{
		colors: make(map[string]color.RGBA),
	}

	return &colors

}

func (c *Color) Get(label string) color.RGBA {
	if col, ok := c.colors[label]; ok {
		return col
	}

	cc := newColor(label)
	c.colors[label] = cc
	return c.colors[label]
}

func newColor(label string) color.RGBA {
	b := hash(label)
	return color.RGBA{b[0], b[1], b[2], 255}
}

func hash(s string) []uint8 {
	h := fnv.New32a()
	h.Write([]byte(s))
	x := h.Sum32()
	b := make([]uint8, 8)
	binary.PutVarint(b, int64(x))
	return b
}
