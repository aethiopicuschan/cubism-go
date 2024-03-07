package drawable

type Drawable struct {
	Id              string
	Texture         int32
	VertexPositions []Vector2
	VertexUvs       []Vector2
	VertexIndices   []uint16
	ConstantFlag    ConstantFlag
	DynamicFlag     DynamicFlag
	Opacity         float32
	Masks           []int32
}
