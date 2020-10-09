package math

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLinearScaleElementNoMin(t *testing.T) {

	scale := newMonotonicMapper(
		1000,
		NewV(0, 0),
		NewV(100, 300))

	assertXMapper(t, scale, 500, 50)

	assertYMapper(t, scale, 167, 50)

}

func TestLinearScaleElementWithMin(t *testing.T) {

	scale := newMonotonicMapper(
		1000,
		NewV(900, 100),
		NewV(1000, 300))

	assertXMapper(t, scale, 0, 900)
	assertXMapper(t, scale, 1000, 1000)
	assertXMapper(t, scale, 500, 950)

	assertYMapper(t, scale, 250, 150)
	assertYMapper(t, scale, 330, 166)

}

func assertXMapper(t *testing.T, mapper Mapper, expected, actual float32) {
	x := mapper.ScaleAt(0, Normal)(actual)
	assert.Equal(t, expected, x)
	sx := mapper.DeScaleAt(0, Normal)(x)
	assert.Equal(t, actual, sx)
}

func assertYMapper(t *testing.T, mapper Mapper, expected, actual float32) {
	y := mapper.ScaleAt(1, Inverse)(actual)
	assert.Equal(t, float64(expected), math.Round(float64(y)))
	sy := mapper.DeScaleAt(1, Inverse)(y)
	assert.Equal(t, float32(actual), sy)
}

func TestFloat32_PrecisionLossNoPanic(t *testing.T) {

	i := time.Now().Unix()

	f := Float32(float64(i))

	assert.NotEqual(t, i, f)

}

func TestMustFloat32_Panic(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	i := time.Now().Unix()

	MustFloat32(float64(i))

}

func TestMustFloat32_LowPrecisionNoPanic(t *testing.T) {

	// if we dont divide by 100 code will cause panic
	i := 1.000000000000002

	MustFloat32(i)

}

func TestMustFloat32_NoPanic(t *testing.T) {

	// if we dont divide by 100 code will cause panic
	i := time.Now().Unix() / 100

	MustFloat32(float64(i))

}
