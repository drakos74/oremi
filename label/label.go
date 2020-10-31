package label

import (
	"fmt"
	"time"
)

const secondsInADay = time.Hour * 24

// Label represents an a dimensional attribute.
// name is the name of the dimension.
// format is the format of the values to be represented as.
type Label struct {
	name   string
	format func(value float32) string
}

func (l Label) Name() string {
	return l.name
}

func (l Label) Format() func(value float32) string {
	return l.format
}

func Num(name string) Label {
	return Label{
		name: name,
		format: func(value float32) string {
			return fmt.Sprintf("%.2f", value)
		},
	}
}

func Date(name string) Label {
	return Label{
		name: name,
		format: func(value float32) string {
			return time.Unix(int64(value), 0).Format(time.Stamp)
		},
	}
}

func Day(name string) Label {
	return Label{
		name: name,
		format: func(value float32) string {
			return time.Unix(int64(value*float32(secondsInADay.Seconds())), 0).Format(time.Stamp)
		},
	}
}
