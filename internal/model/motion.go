package model

import "github.com/aethiopicuschan/cubism-go/internal/motion"

type Meta struct {
	Duration             float64 `json:"Duration"`
	Loop                 bool    `json:"Loop"`
	AreBeziersRestricted bool    `json:"AreBeziersRestricted"`
}

// *.motion3.json用の構造体
type MotionJson struct {
	Version  int  `json:"Version"`
	Meta     Meta `json:"Meta"`
	UserData []struct {
		Time  float64 `json:"Time"`
		Value string  `json:"Value"`
	} `json:"UserData"`
	Curves []struct {
		Target      string    `json:"Target"`
		Id          string    `json:"Id"`
		FadeInTime  *float64  `json:"FadeInTime"`
		FadeOutTime *float64  `json:"FadeOutTime"`
		Segments    []float64 `json:"Segments"`
	} `json:"Curves"`
}

// motion3.jsonをMotionに変換する
func (m *MotionJson) ToMotion(fp string, fadein, fadeout float64, sound string) (mtn motion.Motion) {
	mtn = motion.Motion{
		File:        fp,
		FadeInTime:  fadein,
		FadeOutTime: fadeout,
		Sound:       sound,
		Meta: motion.Meta{
			Duration:             m.Meta.Duration,
			Loop:                 m.Meta.Loop,
			AreBeziersRestricted: m.Meta.AreBeziersRestricted,
		},
	}
	for _, curve := range m.Curves {
		var c motion.Curve
		c.Target = curve.Target
		c.Id = curve.Id
		if curve.FadeInTime == nil {
			c.FadeInTime = -1.0
		} else {
			c.FadeInTime = *curve.FadeInTime
		}
		if curve.FadeOutTime == nil {
			c.FadeOutTime = -1.0
		} else {
			c.FadeOutTime = *curve.FadeOutTime
		}
		var lastPoint motion.Point
		for i := 0; i < len(curve.Segments); {
			if i == 0 {
				lastPoint = motion.Point{
					Time:  curve.Segments[i],
					Value: curve.Segments[i+1],
				}
				i += 2
			}
			segment := curve.Segments[i]
			switch segment {
			case motion.Linear:
				nextPoint := motion.Point{
					Time:  curve.Segments[i+1],
					Value: curve.Segments[i+2],
				}
				c.Segments = append(c.Segments, motion.Segment{
					Points: []motion.Point{
						lastPoint,
						nextPoint,
					},
					Type: motion.Linear,
				})
				lastPoint = nextPoint
				i += 3
			case motion.Bezier:
				t0 := curve.Segments[i+1]
				v0 := curve.Segments[i+2]
				t1 := curve.Segments[i+3]
				v1 := curve.Segments[i+4]
				t2 := curve.Segments[i+5]
				v2 := curve.Segments[i+6]

				nextPoint := motion.Point{
					Time:  t2,
					Value: v2,
				}

				c.Segments = append(c.Segments, motion.Segment{
					Points: []motion.Point{
						lastPoint,
						{
							Time:  t0,
							Value: v0,
						},
						{
							Time:  t1,
							Value: v1,
						},
						nextPoint,
					},
					Type: motion.Bezier,
				})
				lastPoint = nextPoint
				i += 7
			case motion.Stepped:
				t0 := curve.Segments[i+1]
				v0 := curve.Segments[i+2]

				c.Segments = append(c.Segments, motion.Segment{
					Points: []motion.Point{
						lastPoint,
					},
					Type:  motion.Stepped,
					Value: t0,
				})
				lastPoint = motion.Point{
					Time:  t0,
					Value: v0,
				}
				i += 3
			case motion.InverseStepped:
				t0 := curve.Segments[i+1]
				v0 := curve.Segments[i+2]
				tn := lastPoint.Time
				lastPoint = motion.Point{
					Time:  t0,
					Value: v0,
				}
				c.Segments = append(c.Segments, motion.Segment{
					Points: []motion.Point{
						lastPoint,
					},
					Type:  motion.InverseStepped,
					Value: tn,
				})
				i += 3
			}
		}
		mtn.Curves = append(mtn.Curves, c)
	}
	return
}
