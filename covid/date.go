package covid

import (
	"math"
	"time"
)

const secondsInADay = time.Hour * 24

// toDay converts the time to a day based bucket
func toDay(time time.Time) float64 {
	return math.Floor(float64(time.Unix()) / secondsInADay.Seconds())
}
