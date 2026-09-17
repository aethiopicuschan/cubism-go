package core_test

import (
	"testing"

	"github.com/aethiopicuschan/cubism-go/internal/core/core_5"
	"github.com/aethiopicuschan/cubism-go/internal/core/drawable"
	"github.com/stretchr/testify/require"
)

func TestParameterAndPartPointers(t *testing.T) {
	first, second := []byte("first\x00"), []byte("second\x00")
	ids := []*byte{&first[0], &second[0]}
	mins, maxs := []float32{-1, -2}, []float32{1, 2}
	defaults, values := []float32{0, 1}, []float32{0.5, 1.5}
	opacities := []float32{1, 0.5}
	c := core.NewTestCore(core.TestBindings{
		GetParameterCount:         func(uintptr) int { return len(ids) },
		GetParameterIds:           func(uintptr) **byte { return &ids[0] },
		GetParameterMinimumValues: func(uintptr) *float32 { return &mins[0] },
		GetParameterMaximumValues: func(uintptr) *float32 { return &maxs[0] },
		GetParameterDefaultValues: func(uintptr) *float32 { return &defaults[0] },
		GetParameterValues:        func(uintptr) *float32 { return &values[0] },
		GetPartCount:              func(uintptr) int { return len(ids) },
		GetPartIds:                func(uintptr) **byte { return &ids[0] },
		GetPartOpacities:          func(uintptr) *float32 { return &opacities[0] },
	})
	parameters := c.GetParameters(0)
	require.Len(t, parameters, 2)
	require.Equal(t, "second", parameters[1].Id)
	require.Equal(t, mins[1], parameters[1].Minimum)
	require.Equal(t, maxs[1], parameters[1].Maximum)
	require.Equal(t, defaults[1], parameters[1].Default)
	require.Equal(t, values[1], parameters[1].Current)
	require.Equal(t, float32(1.5), c.GetParameterValue(0, "second"))
	c.SetParameterValue(0, "second", 0.75)
	require.Equal(t, []float32{0.5, 0.75}, values)
	require.Zero(t, c.GetParameterValue(0, "missing"))
	c.SetParameterValue(0, "missing", 9)
	require.Equal(t, []float32{0.5, 0.75}, values)
	require.Equal(t, []string{"first", "second"}, c.GetPartIds(0))
	c.SetPartOpacity(0, "second", 0.25)
	c.SetPartOpacity(0, "missing", 9)
	require.Equal(t, []float32{1, 0.25}, opacities)
}

func TestDrawablePointers(t *testing.T) {
	first, second := []byte("first\x00"), []byte("second\x00")
	ids := []*byte{&first[0], &second[0]}
	flags := []uint8{0, 1}
	textures, orders := []int32{3, 7}, []int32{1, 0}
	opacities := []float32{1, 0.5}
	counts := []int32{0, 2}
	positions := []drawable.Vector2{{X: 1, Y: 2}, {X: 3, Y: 4}}
	uvs := []drawable.Vector2{{X: 0, Y: 0}, {X: 1, Y: 1}}
	positionPtrs := []*drawable.Vector2{nil, &positions[0]}
	uvPtrs := []*drawable.Vector2{nil, &uvs[0]}
	indices := []uint16{1, 0}
	indexPtrs := []*uint16{nil, &indices[0]}
	maskCounts := []int32{0, 1}
	mask := int32(0)
	maskPtrs := []*int32{nil, &mask}
	c := core.NewTestCore(core.TestBindings{
		GetDrawableCount:           func(uintptr) int { return 2 },
		GetDrawableIds:             func(uintptr) **byte { return &ids[0] },
		GetDrawableConstantFlags:   func(uintptr) *uint8 { return &flags[0] },
		GetDrawableDynamicFlags:    func(uintptr) *uint8 { return &flags[0] },
		GetDrawableTextureIndices:  func(uintptr) *int32 { return &textures[0] },
		GetDrawableRenderOrders:    func(uintptr) *int32 { return &orders[0] },
		GetDrawableOpacities:       func(uintptr) *float32 { return &opacities[0] },
		GetDrawableVertexCounts:    func(uintptr) *int32 { return &counts[0] },
		GetDrawableVertexPositions: func(uintptr) **drawable.Vector2 { return &positionPtrs[0] },
		GetDrawableVertexUvs:       func(uintptr) **drawable.Vector2 { return &uvPtrs[0] },
		GetDrawableIndexCounts:     func(uintptr) *int32 { return &counts[0] },
		GetDrawableIndices:         func(uintptr) **uint16 { return &indexPtrs[0] },
		GetDrawableMaskCounts:      func(uintptr) *int32 { return &maskCounts[0] },
		GetDrawableMasks:           func(uintptr) **int32 { return &maskPtrs[0] },
	})
	ds := c.GetDrawables(0)
	require.Len(t, ds, 2)
	require.Empty(t, ds[0].VertexPositions)
	require.Empty(t, ds[0].VertexIndices)
	require.Empty(t, ds[0].Masks)
	require.Equal(t, "second", ds[1].Id)
	require.Equal(t, textures[1], ds[1].Texture)
	require.Equal(t, positions, ds[1].VertexPositions)
	require.Equal(t, uvs, ds[1].VertexUvs)
	require.Equal(t, indices, ds[1].VertexIndices)
	require.Equal(t, []int32{0}, ds[1].Masks)
	require.Equal(t, drawable.ParseConstantFlag(flags[1]), ds[1].ConstantFlag)
	require.Equal(t, drawable.ParseDynamicFlag(flags[1]), ds[1].DynamicFlag)
	require.Equal(t, opacities[1], ds[1].Opacity)
	require.Equal(t, positions, c.GetVertexPositions(0)[1])
	require.Equal(t, []int{1, 0}, c.GetSortedDrawableIndices(0))
	positions[0].X = 9
	require.Equal(t, float32(9), ds[1].VertexPositions[0].X)
}

func TestEmptyPointerArrays(t *testing.T) {
	c := core.NewTestCore(core.TestBindings{
		GetParameterCount:          func(uintptr) int { return 0 },
		GetParameterIds:            func(uintptr) **byte { return nil },
		GetParameterValues:         func(uintptr) *float32 { return nil },
		GetDrawableCount:           func(uintptr) int { return 0 },
		GetDrawableVertexCounts:    func(uintptr) *int32 { return nil },
		GetDrawableVertexPositions: func(uintptr) **drawable.Vector2 { return nil },
	})
	require.Empty(t, c.GetVertexPositions(0))
	require.Zero(t, c.GetParameterValue(0, "missing"))
	c.SetParameterValue(0, "missing", 1)
}
