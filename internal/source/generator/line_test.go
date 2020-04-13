package generator

import (
	"fmt"
	"github/drakos74/oremi/internal/model"
	"testing"
)

func TestNewLine(t *testing.T) {

	line := model.NewSeries(4)

	linear := LinearGenerator(1, 1, 2, 4, 2)

	start := model.NewPoint(0, 0, 0, 0)

	line.Add(start)

	for i := 0; i < 100; i++ {

		n, _ := line.Next()

		line.Add(linear.Next(n))

	}

	println(fmt.Sprintf("line := %v", line))

}
