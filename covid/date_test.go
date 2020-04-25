package covid

import (
	"fmt"
	"testing"
	"time"
)

func TestToDay(t *testing.T) {

	d := toDay(time.Now())
	println(fmt.Sprintf("d = %v", d))

}
