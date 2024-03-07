package cubism

import (
	"fmt"

	"github.com/aethiopicuschan/cubism-go/internal/core"
	"github.com/aethiopicuschan/cubism-go/internal/core/drawable"
	"github.com/aethiopicuschan/cubism-go/internal/core/moc"
	"github.com/aethiopicuschan/cubism-go/internal/core/parameter"
	"github.com/aethiopicuschan/cubism-go/internal/motion"
)

// モデルを表す構造体
type Model struct {
	// 内部的に必要なもの
	motionManager *motion.MotionManager
	// 外に見えていいもの
	Version int
	Opacity float32
	// Getterで取得のみ可能なもの
	core          core.Core
	moc           moc.Moc
	textures      []string
	motions       map[string][]*motion.Motion
	sortedIndices []int
	drawables     []Drawable
	// 外に見せるか未定のもの
	physics  physicsJson
	pose     poseJson
	cdi      cdiJson
	exps     []expJson
	userdata userdataJson
}

// コアを取得する
func (m *Model) GetCore() core.Core {
	return m.core
}

// Mocを取得する
func (m *Model) GetMoc() moc.Moc {
	return m.moc
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
	for _, drawable := range m.drawables {
		if drawable.Id == id {
			d = drawable
			return
		}
	}
	err = fmt.Errorf("Drawable not found: %s", id)
	return
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
func (m *Model) GetMotions(groupName string) []*motion.Motion {
	return m.motions[groupName]
}

// モーションを再生する
func (m *Model) PlayMotion(groupName string, index int) {
	m.motionManager = motion.NewMotionManager(m.core, m.moc.ModelPtr, m.motions[groupName][index], func() {
		m.motionManager = nil
	})
}

// モデルを更新する
func (m *Model) Update(delta float64) {
	if m.motionManager != nil {
		m.motionManager.Update(delta)
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
