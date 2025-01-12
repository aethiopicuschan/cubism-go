package motion

import (
	"github.com/aethiopicuschan/cubism-go/internal/core"
)

type MotionManager struct {
	core            core.Core
	modelPtr        uintptr
	queue           []Entry
	lastId          int
	onFinished      func(int)
	savedParameters map[string]float32
}

func NewMotionManager(core core.Core, modelPtr uintptr, onFinished func(int)) *MotionManager {
	return &MotionManager{
		core:            core,
		modelPtr:        modelPtr,
		queue:           []Entry{},
		lastId:          0,
		onFinished:      onFinished,
		savedParameters: make(map[string]float32),
	}
}

func (mm *MotionManager) Start(motion Motion) int {
	mm.lastId++
	mm.queue = append(mm.queue, Entry{
		motion:      motion,
		id:          mm.lastId,
		currentTime: 0,
	})
	return mm.lastId
}

func (mm *MotionManager) Close(id int) {
	index := -1
	for i, entry := range mm.queue {
		if entry.id == id {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	mm.queue[index].motion.LoadedSound.Close()
	mm.queue = append(mm.queue[:index], mm.queue[index+1:]...)
}

func (mm *MotionManager) Reset(id int) {
	index := -1
	for i, entry := range mm.queue {
		if entry.id == id {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	mm.queue[index].currentTime = 0
}

func (mm *MotionManager) saveParameters() {
	parameters := mm.core.GetParameters(mm.modelPtr)
	savedParameters := make(map[string]float32)
	for _, parameter := range parameters {
		savedParameters[parameter.Id] = parameter.Current
	}
	mm.savedParameters = savedParameters
}

func (mm *MotionManager) loadParameters() {
	if mm.savedParameters == nil {
		return
	}
	for id, value := range mm.savedParameters {
		mm.core.SetParameterValue(mm.modelPtr, id, value)
	}
}

func (mm *MotionManager) Update(deltaTime float64) {
	if len(mm.queue) == 0 {
		return
	}
	finished := mm.queue[len(mm.queue)-1].Update(deltaTime)
	if finished {
		mm.onFinished(mm.queue[len(mm.queue)-1].id)
	}
	if len(mm.queue) == 0 {
		return
	}
	mm.loadParameters()

	entry := mm.queue[len(mm.queue)-1]
	if entry.currentTime == deltaTime {
		if entry.motion.Sound != "" {
			entry.motion.LoadedSound.Play()
		}
	}
	fadeIn, fadeOut, fadeWeight := getFade(entry.motion, 1.0, entry.currentTime)
	for _, curve := range entry.motion.Curves {
		for _, seg := range curve.Segments {
			if !segmentIntersects(seg, entry.currentTime) {
				continue
			}
			value := segmentInterpolate(seg, entry.currentTime)
			if curve.Target == "Model" {
				// TODO implement
			}
			if curve.Target == "PartOpacity" {
				mm.core.SetPartOpacity(mm.modelPtr, curve.Id, float32(value))
			}
			if curve.Target == "Parameter" {
				var v float32
				sourceValue := mm.core.GetParameterValue(mm.modelPtr, curve.Id)
				if curve.FadeInTime < 0.0 && curve.FadeOutTime < 0.0 {
					// If the fade is not set for the parameter, apply the motion fade
					v = sourceValue + (float32(value)-sourceValue)*float32(fadeWeight)
				} else {
					// If a fade is set for the parameter, apply that fade
					var fin, fout float64
					if curve.FadeInTime < 0 {
						fin = fadeIn
					} else {
						if curve.FadeInTime == 0.0 {
							fin = 1.0
						} else {
							fin = getEasingSine(entry.currentTime / curve.FadeInTime)
						}
					}
					if curve.FadeOutTime < 0 {
						fout = fadeOut
					} else {
						if curve.FadeOutTime == 0.0 {
							fout = 1.0
						} else {
							fout = getEasingSine((entry.motion.Meta.Duration - entry.currentTime) / curve.FadeOutTime)
						}
					}
					paramWeight := 1.0 * fin * fout
					v = sourceValue + (float32(value)-sourceValue)*float32(paramWeight)
				}
				mm.core.SetParameterValue(mm.modelPtr, curve.Id, v)
			}
		}
	}
	mm.saveParameters()
}
