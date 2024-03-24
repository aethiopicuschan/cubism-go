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

func lerpPoints(a Point, b Point, t float64) Point {
	return Point{
		Time:  a.Time + (b.Time-a.Time)*t,
		Value: a.Value + (b.Value-a.Value)*t,
	}
}

func segmentIntersects(segment Segment, t float64) bool {
	if segment.Type == Linear {
		return segment.Points[0].Time <= t && t <= segment.Points[1].Time
	}
	if segment.Type == Bezier {
		return segment.Points[0].Time <= t && t <= segment.Points[3].Time
	}
	if segment.Type == Stepped {
		return segment.Points[0].Time <= t && t <= segment.Value
	}
	if segment.Type == InverseStepped {
		return segment.Value <= t && t <= segment.Points[0].Time
	}
	return false
}

func segmentInterpolate(segment Segment, t float64) float64 {
	if segment.Type == Linear {
		p0, p1 := segment.Points[0], segment.Points[1]
		k := (t - p0.Time) / (p1.Time - p0.Time)
		if k < 0.0 {
			k = 0.0
		}
		return p0.Value + (p1.Value-p0.Value)*k
	}
	if segment.Type == Bezier {
		p0, p1, p2, p3 := segment.Points[0], segment.Points[1], segment.Points[2], segment.Points[3]
		k := (t - p0.Time) / (p3.Time - p0.Time)
		if k < 0.0 {
			k = 0.0
		}

		p01 := lerpPoints(p0, p1, k)
		p12 := lerpPoints(p1, p2, k)
		p23 := lerpPoints(p2, p3, k)
		p012 := lerpPoints(p01, p12, k)
		p123 := lerpPoints(p12, p23, k)
		return lerpPoints(p012, p123, k).Value
	}
	if segment.Type == Stepped {
		return segment.Points[0].Value
	}
	if segment.Type == InverseStepped {
		return segment.Points[0].Value
	}
	return 0
}

func getFade(motion Motion, weight float64, t float64) (fadeIn, fadeOut, fadeWeight float64) {
	fadeWeight = weight
	if motion.FadeInTime == 0.0 {
		fadeIn = 1.0
	} else {
		fadeIn = getEasingSine(t / motion.FadeInTime)
	}
	if motion.FadeOutTime == 0.0 || motion.Meta.Duration < 0.0 {
		fadeOut = 1.0
	} else {
		fadeOut = getEasingSine((motion.Meta.Duration - t) / motion.FadeOutTime)
	}
	fadeWeight = fadeWeight * fadeIn * fadeOut
	return
}
