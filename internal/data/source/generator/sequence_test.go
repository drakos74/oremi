package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLinearIterator(t *testing.T) {

	linearIterator := NewLinearSequence(0, 0.9)

	v, ok, next := linearIterator.Next()

	assert.True(t, ok)

	assert.Equal(t, 1, v.Dim())

	assert.True(t, next)

}
