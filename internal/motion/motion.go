package motion

const (
	Linear = iota
	Bezier
	Stepped
	InverseStepped
)

type Meta struct {
	Duration             float64
	Loop                 bool
	AreBeziersRestricted bool
}

type Point struct {
	Time  float64
	Value float64
}

type Segment struct {
	Points []Point
	Type   int
	Value  float64
}

type Curve struct {
	Target      string
	Id          string
	FadeInTime  float64
	FadeOutTime float64
	Segments    []Segment
}

type Motion struct {
	File        string
	FadeInTime  float64
	FadeOutTime float64
	Sound       string
	Meta        Meta
	Curves      []Curve
}
