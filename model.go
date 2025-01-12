package cubism

import (
	"fmt"

	"github.com/aethiopicuschan/cubism-go/internal/blink"
	"github.com/aethiopicuschan/cubism-go/internal/core"
	"github.com/aethiopicuschan/cubism-go/internal/core/drawable"
	"github.com/aethiopicuschan/cubism-go/internal/core/moc"
	"github.com/aethiopicuschan/cubism-go/internal/core/parameter"
	"github.com/aethiopicuschan/cubism-go/internal/model"
	"github.com/aethiopicuschan/cubism-go/internal/motion"
)

// A model struct
type Model struct {
	// Internally required
	motionManager *motion.MotionManager
	loopMotions   []int
	blinkManager  *blink.BlinkManager
	// Read-only via getters
	version       int
	core          core.Core
	moc           moc.Moc
	opacity       float32
	textures      []string
	motions       map[string][]motion.Motion
	sortedIndices []int
	drawables     []Drawable
	drawablesMap  map[string]Drawable
	hitAreas      []model.HitArea
	// Not exposed externally
	groups   []model.Group
	physics  model.PhysicsJson
	pose     model.PoseJson
	cdi      model.CdiJson
	exps     []model.ExpJson
	userdata model.UserDataJson
}

// Get the version of the model
func (m *Model) GetVersion() int {
	return m.version
}

// Get the core
func (m *Model) GetCore() core.Core {
	return m.core
}

// Get the moc
func (m *Model) GetMoc() moc.Moc {
	return m.moc
}

// Get the opacity of the model
func (m *Model) GetOpacity() float32 {
	return m.opacity
}

// Get the path of a texture image
func (m *Model) GetTextures() []string {
	return m.textures
}

// Get the sorted drawing order indices
func (m *Model) GetSortedIndices() []int {
	return m.sortedIndices
}

// Get the drawables
func (m *Model) GetDrawables() []Drawable {
	return m.drawables
}

// Get the Drawable with the specified ID
func (m *Model) GetDrawable(id string) (d Drawable, err error) {
	if d, ok := m.drawablesMap[id]; ok {
		return d, nil
	}
	err = fmt.Errorf("Drawable not found: %s", id)
	return
}

// Get the list of hit areas
func (m *Model) GetHitAreas() []model.HitArea {
	return m.hitAreas
}

// Get the list of parameters
func (m *Model) GetParameters() []parameter.Parameter {
	return m.core.GetParameters(m.moc.ModelPtr)
}

// Get the value of the parameter
func (m *Model) GetParameterValue(id string) float32 {
	return m.core.GetParameterValue(m.moc.ModelPtr, id)
}

// Set the value of the parameter
func (m *Model) SetParameterValue(id string, value float32) {
	m.core.SetParameterValue(m.moc.ModelPtr, id, value)
}

// Get the list of motion group names
func (m *Model) GetMotionGroupNames() (names []string) {
	for k := range m.motions {
		names = append(names, k)
	}
	return
}

// Get the list of motions in the group
func (m *Model) GetMotions(groupName string) []motion.Motion {
	return m.motions[groupName]
}

// Play a motion
func (m *Model) PlayMotion(groupName string, index int, loop bool) (id int) {
	if m.motionManager == nil {
		m.motionManager = motion.NewMotionManager(m.core, m.moc.ModelPtr, func(id int) {
			for _, loopId := range m.loopMotions {
				if id == loopId {
					m.motionManager.Reset(id)
					return
				}
			}
			m.motionManager.Close(id)
		})
	}
	id = m.motionManager.Start(m.motions[groupName][index])
	if loop {
		m.loopMotions = append(m.loopMotions, id)
	}
	return
}

// Stop a motion
func (m *Model) StopMotion(id int) {
	for i, loopId := range m.loopMotions {
		if id == loopId {
			m.loopMotions = append(m.loopMotions[:i], m.loopMotions[i+1:]...)
			break
		}
	}
	m.motionManager.Close(id)
}

// Enable Auto Blink
func (m *Model) EnableAutoBlink() {
	for _, group := range m.groups {
		if group.Name == "EyeBlink" {
			m.blinkManager = blink.NewBlinkManager(m.core, m.moc.ModelPtr, group.Ids)
			return
		}
	}
}

// Disable Auto Blink
func (m *Model) DisableAutoBlink() {
	m.blinkManager = nil
}

// Update the model
func (m *Model) Update(delta float64) {
	if m.motionManager != nil {
		m.motionManager.Update(delta)
	}
	if m.blinkManager != nil {
		m.blinkManager.Update(delta)
	}
	m.core.Update(m.moc.ModelPtr)

	// Get the updated dynamic flags
	dfs := m.core.GetDynamicFlags(m.moc.ModelPtr)

	// Check the flags to confirm if there are any targets that need to be updated.
	drawOrderDidChange := false
	renderOrderDidChange := false
	opacityDidChange := false
	vertexPositionsDidChange := false
	// blendColorDidChange := false
	for i := range m.drawables {
		m.drawables[i].DynamicFlag = dfs[i]
		if dfs[i].DrawOrderDidChange {
			drawOrderDidChange = true
		}
		if dfs[i].RenderOrderDidChange {
			renderOrderDidChange = true
		}
		if dfs[i].OpacityDidChange {
			opacityDidChange = true
		}
		if dfs[i].VertexPositionsDidChange {
			vertexPositionsDidChange = true
		}
		/*
			if dfs[i].BlendColorDidChange {
				blendColorDidChange = true
			}
		*/
	}

	// Update the drawing order
	if drawOrderDidChange || renderOrderDidChange {
		m.sortedIndices = m.core.GetSortedDrawableIndices(m.moc.ModelPtr)
	}
	// Update the opacities
	var opacities []float32
	if opacityDidChange {
		opacities = m.core.GetOpacities(m.moc.ModelPtr)
	}
	// Update the vertex positions
	var vertexPositions [][]drawable.Vector2
	if vertexPositionsDidChange {
		vertexPositions = m.core.GetVertexPositions(m.moc.ModelPtr)
	}
	// Update the multiplication color and screen color
	// TODO impl

	for i := range m.drawables {
		if opacityDidChange {
			m.drawables[i].Opacity = opacities[i]
		}
		if vertexPositionsDidChange {
			m.drawables[i].VertexPositions = vertexPositions[i]
		}
	}
}
