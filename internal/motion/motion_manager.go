package motion

import (
	"github.com/aethiopicuschan/cubism-go/internal/core"
	"github.com/aethiopicuschan/cubism-go/internal/sound"
)

type MotionManager struct {
	core            core.Core
	motion          *Motion
	modelPtr        uintptr
	currentTime     float64
	weight          float64
	finished        bool
	onFinished      func()
	savedParameters map[string]float32
}

func NewMotionManager(core core.Core, modelPtr uintptr, motion *Motion, onFinished func()) *MotionManager {
	return &MotionManager{
		core:            core,
		motion:          motion,
		modelPtr:        modelPtr,
		currentTime:     0,
		weight:          1,
		finished:        false,
		onFinished:      onFinished,
		savedParameters: make(map[string]float32),
	}
}

func (m *MotionManager) lerpPoints(a Point, b Point, t float64) Point {
	return Point{
		Time:  a.Time + (b.Time-a.Time)*t,
		Value: a.Value + (b.Value-a.Value)*t,
	}
}

func (m *MotionManager) segmentIntersects(segment Segment) bool {
	if segment.Type == Linear {
		return segment.Points[0].Time <= m.currentTime && m.currentTime <= segment.Points[1].Time
	}
	if segment.Type == Bezier {
		return segment.Points[0].Time <= m.currentTime && m.currentTime <= segment.Points[3].Time
	}
	if segment.Type == Stepped {
		return segment.Points[0].Time <= m.currentTime && m.currentTime <= segment.Value
	}
	if segment.Type == InverseStepped {
		return segment.Value <= m.currentTime && m.currentTime <= segment.Points[0].Time
	}
	return false
}

func (m *MotionManager) segmentInterpolate(segment Segment) float64 {
	t := m.currentTime
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

		p01 := m.lerpPoints(p0, p1, k)
		p12 := m.lerpPoints(p1, p2, k)
		p23 := m.lerpPoints(p2, p3, k)
		p012 := m.lerpPoints(p01, p12, k)
		p123 := m.lerpPoints(p12, p23, k)
		return m.lerpPoints(p012, p123, k).Value
	}
	if segment.Type == Stepped {
		return segment.Points[0].Value
	}
	if segment.Type == InverseStepped {
		return segment.Points[0].Value
	}
	return 0
}

func (m *MotionManager) getFade() (fadeIn, fadeOut, fadeWeight float64) {
	fadeWeight = m.weight
	if m.motion.FadeInTime == 0.0 {
		fadeIn = 1.0
	} else {
		fadeIn = getEasingSine(m.currentTime / m.motion.FadeInTime)
	}
	if m.motion.FadeOutTime == 0.0 || m.motion.Meta.Duration < 0.0 {
		fadeOut = 1.0
	} else {
		fadeOut = getEasingSine((m.motion.Meta.Duration - m.currentTime) / m.motion.FadeOutTime)
	}
	fadeWeight = fadeWeight * fadeIn * fadeOut
	return
}

func (m *MotionManager) saveParameters() {
	parameters := m.core.GetParameters(m.modelPtr)
	savedParameters := make(map[string]float32)
	for _, parameter := range parameters {
		savedParameters[parameter.Id] = parameter.Current
	}
	m.savedParameters = savedParameters
}

func (m *MotionManager) loadParameters() {
	if m.savedParameters == nil {
		return
	}
	for id, value := range m.savedParameters {
		m.core.SetParameterValue(m.modelPtr, id, value)
	}
}

func (m *MotionManager) Update(delta float64) (err error) {
	if m.finished {
		return
	}
	if m.currentTime == 0.0 {
		if m.motion.Sound != "" {
			format, err := sound.DetectFormat(m.motion.Sound)
			if err != nil {
				return err
			}
			if err := sound.Play(format, m.motion.LoadedSound); err != nil {
				return err
			}
		}
	}
	m.currentTime += delta
	if m.currentTime >= m.motion.Meta.Duration {
		if m.motion.Meta.Loop {
			m.currentTime = 0.0
		} else {
			m.finished = true
			m.onFinished()
		}
	}

	fadeIn, fadeOut, fadeWeight := m.getFade()

	m.loadParameters()
	for _, curve := range m.motion.Curves {
		for _, seg := range curve.Segments {
			if !m.segmentIntersects(seg) {
				continue
			}
			value := m.segmentInterpolate(seg)
			if curve.Target == "Model" {
				// TODO implement
			}
			if curve.Target == "PartOpacity" {
				m.core.SetPartOpacity(m.modelPtr, curve.Id, float32(value))
			}
			if curve.Target == "Parameter" {
				var v float32
				sourceValue := m.core.GetParameterValue(m.modelPtr, curve.Id)
				if curve.FadeInTime < 0.0 && curve.FadeOutTime < 0.0 {
					// パラメータに対してフェードが設定されていない場合はモーションのフェードを適用する
					v = sourceValue + (float32(value)-sourceValue)*float32(fadeWeight)
				} else {
					// パラメータに対してフェードが設定されている場合はそちらを適用する
					var fin, fout float64
					if curve.FadeInTime < 0 {
						fin = fadeIn
					} else {
						if curve.FadeInTime == 0.0 {
							fin = 1.0
						} else {
							fin = getEasingSine(m.currentTime / curve.FadeInTime)
						}
					}
					if curve.FadeOutTime < 0 {
						fout = fadeOut
					} else {
						if curve.FadeOutTime == 0.0 {
							fout = 1.0
						} else {
							fout = getEasingSine((m.motion.Meta.Duration - m.currentTime) / curve.FadeOutTime)
						}
					}
					paramWeight := m.weight * fin * fout
					v = sourceValue + (float32(value)-sourceValue)*float32(paramWeight)
				}
				m.core.SetParameterValue(m.modelPtr, curve.Id, v)
			}
		}
	}
	m.saveParameters()
	return
}
