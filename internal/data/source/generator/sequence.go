package generator

import (
	"fmt"
	"github/drakos74/oremi/internal/data/model"
)

type Sequence struct {
	f       Y
	count   int
	start   float64
	current float64
	limit   float64
}

func (i *Sequence) Next() (vector model.Vector, ok, hasNext bool) {
	var next float64
	if i.count == 0 {
		next = i.start
	} else {
		next = i.f(i.current)
	}
	i.count++
	i.current = next
	return model.NewVector([]string{fmt.Sprintf("%d", i.count)}, i.current), true, i.limit == 0 || i.current < i.limit
}

func (i *Sequence) Reset() {
	i.current = i.start
	i.count = 0
}

type LinearSequence struct {
	start float64
}

func NewLinearSequence(s float64, step float64) *Sequence {
	return &Sequence{
		f: func(x ...float64) float64 {
			return x[0] + step
		},
		start: s,
	}
}
