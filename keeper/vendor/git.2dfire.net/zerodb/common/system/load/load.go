package load

import (
	"fmt"
)

type Avg struct {
	Avg1min  float64
	Avg5min  float64
	Avg15min float64
}

func (a *Avg) String() string {
	return fmt.Sprintf("<1min:%f, 5min:%f, 15min:%f>", a.Avg1min, a.Avg5min, a.Avg15min)
}
