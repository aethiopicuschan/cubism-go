package core

import "github.com/aethiopicuschan/cubism-go/internal/core/drawable"

type TestBindings struct {
	GetParameterCount          func(uintptr) int
	GetParameterIds            func(uintptr) **byte
	GetParameterMinimumValues  func(uintptr) *float32
	GetParameterMaximumValues  func(uintptr) *float32
	GetParameterDefaultValues  func(uintptr) *float32
	GetParameterValues         func(uintptr) *float32
	GetPartCount               func(uintptr) int
	GetPartIds                 func(uintptr) **byte
	GetPartOpacities           func(uintptr) *float32
	GetDrawableCount           func(uintptr) int
	GetDrawableIds             func(uintptr) **byte
	GetDrawableConstantFlags   func(uintptr) *uint8
	GetDrawableDynamicFlags    func(uintptr) *uint8
	GetDrawableTextureIndices  func(uintptr) *int32
	GetDrawableRenderOrders    func(uintptr) *int32
	GetDrawableOpacities       func(uintptr) *float32
	GetDrawableVertexCounts    func(uintptr) *int32
	GetDrawableVertexPositions func(uintptr) **drawable.Vector2
	GetDrawableVertexUvs       func(uintptr) **drawable.Vector2
	GetDrawableIndexCounts     func(uintptr) *int32
	GetDrawableIndices         func(uintptr) **uint16
	GetDrawableMaskCounts      func(uintptr) *int32
	GetDrawableMasks           func(uintptr) **int32
}

func NewTestCore(b TestBindings) *Core {
	return &Core{
		csmGetParameterCount:          b.GetParameterCount,
		csmGetParameterIds:            b.GetParameterIds,
		csmGetParameterMinimumValues:  b.GetParameterMinimumValues,
		csmGetParameterMaximumValues:  b.GetParameterMaximumValues,
		csmGetParameterDefaultValues:  b.GetParameterDefaultValues,
		csmGetParameterValues:         b.GetParameterValues,
		csmGetPartCount:               b.GetPartCount,
		csmGetPartIds:                 b.GetPartIds,
		csmGetPartOpacities:           b.GetPartOpacities,
		csmGetDrawableCount:           b.GetDrawableCount,
		csmGetDrawableIds:             b.GetDrawableIds,
		csmGetDrawableConstantFlags:   b.GetDrawableConstantFlags,
		csmGetDrawableDynamicFlags:    b.GetDrawableDynamicFlags,
		csmGetDrawableTextureIndices:  b.GetDrawableTextureIndices,
		csmGetDrawableRenderOrders:    b.GetDrawableRenderOrders,
		csmGetDrawableOpacities:       b.GetDrawableOpacities,
		csmGetDrawableVertexCounts:    b.GetDrawableVertexCounts,
		csmGetDrawableVertexPositions: b.GetDrawableVertexPositions,
		csmGetDrawableVertexUvs:       b.GetDrawableVertexUvs,
		csmGetDrawableIndexCounts:     b.GetDrawableIndexCounts,
		csmGetDrawableIndices:         b.GetDrawableIndices,
		csmGetDrawableMaskCounts:      b.GetDrawableMaskCounts,
		csmGetDrawableMasks:           b.GetDrawableMasks,
	}
}
