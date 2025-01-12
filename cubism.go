package cubism

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/aethiopicuschan/cubism-go/internal/core"
	"github.com/aethiopicuschan/cubism-go/internal/model"
	"github.com/aethiopicuschan/cubism-go/internal/motion"
	"github.com/aethiopicuschan/cubism-go/sound"
	"github.com/aethiopicuschan/cubism-go/sound/disabled"
)

/*
The main body of cubism-go
*/
type Cubism struct {
	core core.Core
	// A function to load audio files
	LoadSound func(fp string) (s sound.Sound, err error)
}

// Constructor for the [Cubism] struct
func NewCubism(lib string) (c Cubism, err error) {
	c.core, err = core.NewCore(lib)
	return
}

// Load a model from model3.json
func (c *Cubism) LoadModel(path string) (m *Model, err error) {
	m = &Model{
		core:    c.core,
		opacity: 1.0,
	}

	// Get the absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	// Get the directory
	dir := filepath.Dir(absPath)

	// Read model3.json
	buf, err := os.ReadFile(absPath)
	if err != nil {
		return
	}
	// Convert to a structure compatible with version 3
	var mj model.ModelJson
	if err = json.Unmarshal(buf, &mj); err != nil {
		return
	}

	// Get the version information
	m.version = mj.Version
	// Convert the path of the texture image to an absolute path
	m.textures = mj.FileReferences.Textures
	for i := range m.textures {
		m.textures[i] = filepath.Join(dir, m.textures[i])
	}

	m.groups = mj.Groups
	m.hitAreas = mj.HitAreas

	// Load the moc3 file
	moc3Path := filepath.Join(dir, mj.FileReferences.Moc)
	m.moc, err = c.core.LoadMoc(moc3Path)
	if err != nil {
		return
	}
	// Get the Drawables
	ds := c.core.GetDrawables(m.moc.ModelPtr)
	for _, d := range ds {
		m.drawables = append(m.drawables, Drawable{
			Id:              d.Id,
			Texture:         m.textures[d.Texture],
			VertexPositions: d.VertexPositions,
			VertexUvs:       d.VertexUvs,
			VertexIndices:   d.VertexIndices,
			ConstantFlag:    d.ConstantFlag,
			DynamicFlag:     d.DynamicFlag,
			Opacity:         d.Opacity,
			Masks:           d.Masks,
		})
	}
	// Create map of Drawables
	m.drawablesMap = map[string]Drawable{}
	for _, d := range m.drawables {
		m.drawablesMap[d.Id] = d
	}
	// Get the sorted indices
	m.sortedIndices = c.core.GetSortedDrawableIndices(m.moc.ModelPtr)

	// Load the physics settings if they exist
	if mj.FileReferences.Physics != "" {
		physicsPath := filepath.Join(dir, mj.FileReferences.Physics)
		buf, err = os.ReadFile(physicsPath)
		if err != nil {
			return
		}
		if err = json.Unmarshal(buf, &m.physics); err != nil {
			return
		}
	}

	// Load the pose settings if they exist
	if mj.FileReferences.Pose != "" {
		posePath := filepath.Join(dir, mj.FileReferences.Pose)
		buf, err = os.ReadFile(posePath)
		if err != nil {
			return
		}
		if err = json.Unmarshal(buf, &m.pose); err != nil {
			return
		}
	}

	// Load the display info settings if they exist
	if mj.FileReferences.DisplayInfo != "" {
		displayInfoPath := filepath.Join(dir, mj.FileReferences.DisplayInfo)
		buf, err = os.ReadFile(displayInfoPath)
		if err != nil {
			return
		}
		if err = json.Unmarshal(buf, &m.cdi); err != nil {
			return
		}
	}

	// Load the expressions
	for _, exp := range mj.FileReferences.Expressions {
		expPath := filepath.Join(dir, exp.File)
		buf, err = os.ReadFile(expPath)
		if err != nil {
			return
		}
		var e model.ExpJson
		if err = json.Unmarshal(buf, &e); err != nil {
			return
		}
		e.Name = exp.Name
		m.exps = append(m.exps, e)
	}

	// Load the motion settings
	m.motions = map[string][]motion.Motion{}
	for name, motions := range mj.FileReferences.Motions {
		m.motions[name] = []motion.Motion{}
		for _, motion := range motions {
			motionPath := filepath.Join(dir, motion.File)
			buf, err = os.ReadFile(motionPath)
			if err != nil {
				return
			}
			var mtnJson model.MotionJson
			if err = json.Unmarshal(buf, &mtnJson); err != nil {
				return
			}
			fp := filepath.Base(motion.File)
			motion := mtnJson.ToMotion(fp, motion.FadeInTime, motion.FadeOutTime, motion.Sound)
			if motion.Sound != "" {
				soundPath := filepath.Join(dir, motion.Sound)
				// If LoadSound is nil, don't play the sound
				if c.LoadSound == nil {
					motion.LoadedSound, err = disabled.LoadSound(soundPath)
				} else {
					motion.LoadedSound, err = c.LoadSound(soundPath)
				}
				if err != nil {
					return
				}
			}
			m.motions[name] = append(m.motions[name], motion)
		}
	}

	// Load user data if it exists
	if mj.FileReferences.UserData != "" {
		userDataPath := filepath.Join(dir, mj.FileReferences.UserData)
		buf, err = os.ReadFile(userDataPath)
		if err != nil {
			return
		}
		if err = json.Unmarshal(buf, &m.userdata); err != nil {
			return
		}
	}

	return
}
