package cubism

import "github.com/aethiopicuschan/cubism-go/internal/motion"

// *.model3.json用の構造体
type modelJson struct {
	Version        int `json:"Version"`
	FileReferences struct {
		Moc         string   `json:"Moc"`
		Textures    []string `json:"Textures"`
		Physics     string   `json:"Physics"`
		Pose        string   `json:"Pose"`
		DisplayInfo string   `json:"DisplayInfo"`
		Expressions []struct {
			Name string `json:"Name"`
			File string `json:"File"`
		} `json:"Expressions"`
		Motions map[string][]struct {
			File        string  `json:"File"`
			FadeInTime  float64 `json:"FadeInTime"`
			FadeOutTime float64 `json:"FadeOutTime"`
			Sound       string  `json:"Sound"`
			MotionSync  string  `json:"MotionSync"`
		} `json:"Motions"`
		UserData string `json:"UserData"`
	} `json:"FileReferences"`
	Groups   []group   `json:"Groups"`
	HitAreas []hitArea `json:"HitAreas"`
}

type group struct {
	Target string   `json:"Target"`
	Name   string   `json:"Name"`
	Ids    []string `json:"Ids"`
}

type hitArea struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

// *.physics3.json用の構造体
type physicsJson struct {
	Version int `json:"Version"`
	Meta    struct {
		PhysicsSettingCount int `json:"PhysicsSettingCount"`
		TotalInputCount     int `json:"TotalInputCount"`
		TotalOutputCount    int `json:"TotalOutputCount"`
		VertexCount         int `json:"VertexCount"`
		EffectiveForces     struct {
			Gravity struct {
				X float64 `json:"X"`
				Y float64 `json:"Y"`
			} `json:"Gravity"`
			Wind struct {
				X float64 `json:"X"`
				Y float64 `json:"Y"`
			} `json:"Wind"`
		} `json:"EffectiveForces"`
		PhysicsDictionary []struct {
			Id   string `json:"Id"`
			Name string `json:"Name"`
		} `json:"PhysicsDictionary"`
	} `json:"Meta"`
	PhysicsSettings []struct {
		Id    string `json:"Id"`
		Input []struct {
			Source struct {
				Target string `json:"Target"`
				Id     string `json:"Id"`
			} `json:"Source"`
			Weight  float64 `json:"Weight"`
			Type    string  `json:"Type"`
			Reflect bool    `json:"Reflect"`
		} `json:"Input"`
		Output []struct {
			Destination struct {
				Target string `json:"Target"`
				Id     string `json:"Id"`
			} `json:"Destination"`
			VertexIndex int     `json:"VertexIndex"`
			Scale       float64 `json:"Scale"`
			Weight      float64 `json:"Weight"`
			Type        string  `json:"Type"`
			Reflect     bool    `json:"Reflect"`
		} `json:"Output"`
		Vertices []struct {
			Position struct {
				X float64 `json:"X"`
				Y float64 `json:"Y"`
			} `json:"Position"`
			Mobility     float64 `json:"Mobility"`
			Delay        float64 `json:"Delay"`
			Acceleration float64 `json:"Acceleration"`
			Radius       float64 `json:"Radius"`
		} `json:"Vertices"`
		Normalization struct {
			Position struct {
				Minimum float64 `json:"Minimum"`
				Default float64 `json:"Default"`
				Maximum float64 `json:"Maximum"`
			} `json:"Position"`
			Angle struct {
				Minimum float64 `json:"Minimum"`
				Default float64 `json:"Default"`
				Maximum float64 `json:"Maximum"`
			} `json:"Angle"`
		} `json:"Normalization"`
	} `json:"PhysicsSettings"`
}

// *.pose3.json用の構造体
type poseJson struct {
	Type       string  `json:"Type"`
	FadeInTime float64 `json:"FadeInTime"`
	Groups     [][]struct {
		Id   string   `json:"Id"`
		Link []string `json:"Link"`
	} `json:"Groups"`
}

// *.cdi3.json用の構造体
type cdiJson struct {
	Version    int `json:"Version"`
	Parameters []struct {
		Id      string `json:"Id"`
		GroupId string `json:"GroupId"`
		Name    string `json:"Name"`
	} `json:"Parameters"`
	ParameterGroups []struct {
		Id      string `json:"Id"`
		GroupId string `json:"GroupId"`
		Name    string `json:"Name"`
	} `json:"ParameterGroups"`
	Parts []struct {
		Id   string `json:"Id"`
		Name string `json:"Name"`
	} `json:"Parts"`
}

// *.exp3.json用の構造体
type expJson struct {
	Name       string
	Type       string `json:"Type"`
	Parameters []struct {
		Id    string  `json:"Id"`
		Value float64 `json:"Value"`
		Blend string  `json:"Blend"`
	} `json:"Parameters"`
}

type meta struct {
	Duration             float64 `json:"Duration"`
	Loop                 bool    `json:"Loop"`
	AreBeziersRestricted bool    `json:"AreBeziersRestricted"`
}

// *.motion3.json用の構造体
type motionJson struct {
	Version  int  `json:"Version"`
	Meta     meta `json:"Meta"`
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
func (m *motionJson) toMotion(fp string, fadein, fadeout float64, sound string) (mtn *motion.Motion) {
	mtn = &motion.Motion{}
	mtn.File = fp
	mtn.FadeInTime = fadein
	mtn.FadeOutTime = fadeout
	mtn.Sound = sound
	mtn.Meta.Duration = m.Meta.Duration
	mtn.Meta.Loop = m.Meta.Loop
	mtn.Meta.AreBeziersRestricted = m.Meta.AreBeziersRestricted
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

// *.userdata3.json用の構造体
type userdataJson struct {
	Version int `json:"Version"`
	Meta    struct {
		UserDataCount     int `json:"UserDataCount"`
		TotalUserDataSize int `json:"TotalUserDataSize"`
	} `json:"Meta"`
	UserData []struct {
		Target string `json:"Target"`
		Id     string `json:"Id"`
		Value  string `json:"Value"`
	} `json:"UserData"`
}
