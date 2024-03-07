package blink

import (
	"math/rand"

	"github.com/aethiopicuschan/cubism-go/internal/core"
)

const (
	EyeStateFirst    = iota ///< 初期状態
	EyeStateInterval        ///< まばたきしていない状態
	EyeStateClosing         ///< まぶたが閉じていく途中の状態
	EyeStateClosed          ///< まぶたが閉じている状態
	EyeStateOpening         ///< まぶたが開いていく途中の状態
)

type BlinkManager struct {
	core             core.Core
	modelPtr         uintptr
	ids              []string
	state            int
	interval         float64
	closing          float64
	opening          float64
	currentTime      float64
	stateStartTime   float64
	nextBlinkingTime float64
}

func NewBlinkManager(core core.Core, modelPtr uintptr, ids []string) *BlinkManager {
	return &BlinkManager{
		core:             core,
		modelPtr:         modelPtr,
		ids:              ids,
		state:            EyeStateFirst,
		interval:         4.0,
		closing:          0.1,
		opening:          0.15,
		currentTime:      0,
		stateStartTime:   0,
		nextBlinkingTime: 0,
	}
}

func (b *BlinkManager) DetermineNextBlinkingTiming() float64 {
	r := rand.Float64()
	return b.currentTime + (r * (2.0*b.interval - 1.0))
}

func (b *BlinkManager) Update(delta float64) {
	b.currentTime += delta

	var value float32

	switch b.state {
	case EyeStateFirst:
		b.state = EyeStateInterval
		b.nextBlinkingTime = b.DetermineNextBlinkingTiming()
		value = 1.0
	case EyeStateInterval:
		if b.currentTime >= b.nextBlinkingTime {
			b.state = EyeStateClosing
			b.stateStartTime = b.currentTime
		}
		value = 1.0
	case EyeStateClosing:
		t := (b.currentTime - b.stateStartTime) / b.closing
		if t >= 1 {
			b.state = EyeStateClosed
			b.stateStartTime = b.currentTime
		}
		value = 1.0 - float32(t)
	case EyeStateClosed:
		t := (b.currentTime - b.stateStartTime) / b.closing
		if t >= 1 {
			b.state = EyeStateOpening
			b.stateStartTime = b.currentTime
		}
		value = 0.0
	case EyeStateOpening:
		t := (b.currentTime - b.stateStartTime) / b.opening
		if t >= 1 {
			t = 1
			b.state = EyeStateInterval
			b.nextBlinkingTime = b.DetermineNextBlinkingTiming()
		}
		value = float32(t)
	}

	for _, id := range b.ids {
		b.core.SetParameterValue(b.modelPtr, id, value)
	}
}
