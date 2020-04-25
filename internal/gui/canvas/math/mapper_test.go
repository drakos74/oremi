package math

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"gioui.org/f32"
)

func TestLinearScaleElementNoMin(t *testing.T) {

	scale := NewLinearMapper(
		1000,
		f32.Point{
			X: 0,
			Y: 0,
		}, f32.Point{
			X: 100,
			Y: 300,
		})

	assertXMapper(t, scale, 500, 50)

	assertYMapper(t, scale, 167, 50)

}

func TestLinearScaleElementWithMin(t *testing.T) {

	scale := NewLinearMapper(
		1000,
		f32.Point{
			X: 900,
			Y: 100,
		}, f32.Point{
			X: 1000,
			Y: 300,
		})

	assertXMapper(t, scale, 0, 900)
	assertXMapper(t, scale, 1000, 1000)
	assertXMapper(t, scale, 500, 950)

	assertYMapper(t, scale, 250, 150)
	assertYMapper(t, scale, 330, 166)

}

func assertXMapper(t *testing.T, mapper Mapper, expected, actual float32) {
	x := mapper.ScaleX()(actual)
	assert.Equal(t, expected, x)
	sx := mapper.DeScaleX()(x)
	assert.Equal(t, actual, sx)
}

func assertYMapper(t *testing.T, mapper Mapper, expected, actual float32) {
	y := mapper.ScaleY()(actual)
	assert.Equal(t, float64(expected), math.Round(float64(y)))
	sy := mapper.DeScaleY()(y)
	assert.Equal(t, float32(actual), sy)
}
