package motion

import "math"

func getEasingSine(value float64) float64 {
	if value < 0.0 {
		return 0.0
	} else if value > 1.0 {
		return 1.0
	}

	return 0.5 - 0.5*math.Cos(value*math.Pi)
}
