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

// モデルを表す構造体
type Model struct {
	// 内部的に必要なもの
	motionManager *motion.MotionManager
	loopMotions   []int
	blinkManager  *blink.BlinkManager
	// Getterで取得のみ可能なもの
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
	// 外に見せるか未定のもの
	groups   []model.Group
	physics  model.PhysicsJson
	pose     model.PoseJson
	cdi      model.CdiJson
	exps     []model.ExpJson
	userdata model.UserDataJson
}

// バージョンを取得する
func (m *Model) GetVersion() int {
	return m.version
}

// コアを取得する
func (m *Model) GetCore() core.Core {
	return m.core
}

// Mocを取得する
func (m *Model) GetMoc() moc.Moc {
	return m.moc
}

// 透明度を取得する
func (m *Model) GetOpacity() float32 {
	return m.opacity
}

// テクスチャ画像のパスを取得する
func (m *Model) GetTextures() []string {
	return m.textures
}

// ソート済みの描画順のインデックスを取得する
func (m *Model) GetSortedIndices() []int {
	return m.sortedIndices
}

// Drawablesを取得する
func (m *Model) GetDrawables() []Drawable {
	return m.drawables
}

// 指定したIDのDrawableを取得する
func (m *Model) GetDrawable(id string) (d Drawable, err error) {
	if d, ok := m.drawablesMap[id]; ok {
		return d, nil
	}
	err = fmt.Errorf("Drawable not found: %s", id)
	return
}

// ヒットエリアの一覧を取得する
func (m *Model) GetHitAreas() []model.HitArea {
	return m.hitAreas
}

// パラメータの一覧を取得する
func (m *Model) GetParameters() []parameter.Parameter {
	return m.core.GetParameters(m.moc.ModelPtr)
}

// パラメータの値を取得する
func (m *Model) GetParameterValue(id string) float32 {
	return m.core.GetParameterValue(m.moc.ModelPtr, id)
}

// パラメータの値を設定する
func (m *Model) SetParameterValue(id string, value float32) {
	m.core.SetParameterValue(m.moc.ModelPtr, id, value)
}

// モーションのグループ名の一覧を取得する
func (m *Model) GetMotionGroupNames() (names []string) {
	for k := range m.motions {
		names = append(names, k)
	}
	return
}

// グループに含まれるモーションの一覧を取得する
func (m *Model) GetMotions(groupName string) []motion.Motion {
	return m.motions[groupName]
}

// モーションを再生する
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

// モーションを停止する
func (m *Model) StopMotion(id int) {
	for i, loopId := range m.loopMotions {
		if id == loopId {
			m.loopMotions = append(m.loopMotions[:i], m.loopMotions[i+1:]...)
			break
		}
	}
	m.motionManager.Close(id)
}

// 自動まばたきを有効にする
func (m *Model) EnableAutoBlink() {
	for _, group := range m.groups {
		if group.Name == "EyeBlink" {
			m.blinkManager = blink.NewBlinkManager(m.core, m.moc.ModelPtr, group.Ids)
			return
		}
	}
}

// 自動まばたきを無効にする
func (m *Model) DisableAutoBlink() {
	m.blinkManager = nil
}

// モデルを更新する
func (m *Model) Update(delta float64) {
	if m.motionManager != nil {
		m.motionManager.Update(delta)
	}
	if m.blinkManager != nil {
		m.blinkManager.Update(delta)
	}
	m.core.Update(m.moc.ModelPtr)

	// 更新された動的フラグを取得
	dfs := m.core.GetDynamicFlags(m.moc.ModelPtr)

	// 更新すべき対象が無いかフラグを舐めて確認する
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

	// 描画順の更新
	if drawOrderDidChange || renderOrderDidChange {
		m.sortedIndices = m.core.GetSortedDrawableIndices(m.moc.ModelPtr)
	}
	// 透明度の更新
	var opacities []float32
	if opacityDidChange {
		opacities = m.core.GetOpacities(m.moc.ModelPtr)
	}
	// 頂点の更新
	var vertexPositions [][]drawable.Vector2
	if vertexPositionsDidChange {
		vertexPositions = m.core.GetVertexPositions(m.moc.ModelPtr)
	}
	// 乗算色・スクリーン色の更新
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
