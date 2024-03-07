package cubism

import "github.com/aethiopicuschan/cubism-go/internal/core/drawable"

type Drawable struct {
	Id              string
	Texture         string
	VertexPositions []drawable.Vector2
	VertexUvs       []drawable.Vector2
	VertexIndices   []uint16
	ConstantFlag    drawable.ConstantFlag
	DynamicFlag     drawable.DynamicFlag
	Opacity         float32
	Masks           []int32
}
