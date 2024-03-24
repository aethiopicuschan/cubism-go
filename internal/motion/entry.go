package motion

type Entry struct {
	motion      Motion
	id          int
	currentTime float64
}

func (e *Entry) Update(deltaTime float64) (finished bool) {
	e.currentTime += deltaTime
	if e.currentTime >= e.motion.Meta.Duration {
		finished = true
		return
	}
	return
}
