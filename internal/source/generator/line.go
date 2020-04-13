package generator

import "github/drakos74/oremi/internal/model"

type LinearSeries struct {
	model.Series
}

func NewLine() *LinearSeries {
	return &LinearSeries{*model.NewSeries(2)}
}

func (l *LinearSeries) Generate(num int, grad ...float64) {

	linear := LinearGenerator(1, grad...)

	start := model.NewPoint("0", make([]float64, len(grad))...)

	l.Add(start)

	for i := 0; i < num; i++ {

		n, _ := l.Next()

		l.Add(linear.Next(n))

	}

	l.Reset()
}

func LinearGenerator(s float64, grad ...float64) Generator {
	f := make([]Y, len(grad))
	for i, g := range grad {
		f[i] = func(g float64) func(x float64) float64 {
			return func(x float64) float64 {
				return x + s*g
			}
		}(g)
	}
	return Generator{f: f}
}
