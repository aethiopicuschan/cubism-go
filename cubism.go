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
cubism-goの本体
*/
type Cubism struct {
	core core.Core
	// 音声ファイルを読み込む関数
	LoadSound func(fp string) (s sound.Sound, err error)
}

func NewCubism(lib string) (c Cubism, err error) {
	c.core, err = core.NewCore(lib)
	return
}

// model3.jsonからモデルを読み込む
func (c *Cubism) LoadModel(path string) (m *Model, err error) {
	m = &Model{
		core:    c.core,
		opacity: 1.0,
	}

	// 絶対パスを取得
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	// ディレクトリを取得
	dir := filepath.Dir(absPath)

	// model3.jsonを読み込む
	buf, err := os.ReadFile(absPath)
	if err != nil {
		return
	}
	// バージョン3に対応した構造体にする
	var mj model.ModelJson
	if err = json.Unmarshal(buf, &mj); err != nil {
		return
	}

	// バージョン情報
	m.version = mj.Version
	// テクスチャ画像のパスを絶対パスにする
	m.textures = mj.FileReferences.Textures
	for i := range m.textures {
		m.textures[i] = filepath.Join(dir, m.textures[i])
	}

	m.groups = mj.Groups
	m.hitAreas = mj.HitAreas

	// moc3ファイルを読み込む
	moc3Path := filepath.Join(dir, mj.FileReferences.Moc)
	m.moc, err = c.core.LoadMoc(moc3Path)
	if err != nil {
		return
	}
	// Drawablesを取得
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
	// Drawablesのmapを作成
	m.drawablesMap = map[string]Drawable{}
	for _, d := range m.drawables {
		m.drawablesMap[d.Id] = d
	}
	// ソート済みインデックスを取得
	m.sortedIndices = c.core.GetSortedDrawableIndices(m.moc.ModelPtr)

	// 物理演算の設定を読み込む(あれば)
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

	// ポーズの設定を読み込む(あれば)
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

	// 表示補助の設定を読み込む(あれば)
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

	// 表情の設定を読み込む
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

	// モーションの設定を読み込む
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
				// もしLoadSoundがnilなら音声は流さない
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

	// ユーザデータを読み込む(あれば)
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
