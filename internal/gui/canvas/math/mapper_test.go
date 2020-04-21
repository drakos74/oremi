package math

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"gioui.org/f32"
)

func TestLinearScaleElement(t *testing.T) {

	scale := NewLinearMapper(
		1000,
		f32.Point{
			X: 0,
			Y: 0,
		}, f32.Point{
			X: 100,
			Y: 300,
		})

	x := scale.ScaleX()(50)
	assert.Equal(t, float32(500), x)
	sx := scale.DeScaleX()(x)
	assert.Equal(t, float32(50), sx)

	y := scale.ScaleY()(50)
	assert.Equal(t, float64(167), math.Round(float64(y)))
	sy := scale.DeScaleY()(y)
	assert.Equal(t, float32(50), sy)

}
