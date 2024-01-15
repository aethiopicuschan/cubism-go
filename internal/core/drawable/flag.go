package drawable

type ConstantFlag struct {
	BlendAdditive       bool
	BlendMultiplicative bool
	IsDoubleSided       bool
	IsInvertedMask      bool
}

func ParseConstantFlag(flag uint8) (c ConstantFlag) {
	c.BlendAdditive = flag&1 == 1
	c.BlendMultiplicative = flag&2 == 2
	c.IsDoubleSided = flag&4 == 4
	c.IsInvertedMask = flag&8 == 8
	return
}

type DynamicFlag struct {
	IsVisible                bool
	VisibilityDidChange      bool
	OpacityDidChange         bool
	DrawOrderDidChange       bool
	RenderOrderDidChange     bool
	VertexPositionsDidChange bool
	BlendColorDidChange      bool
}

func ParseDynamicFlag(flag uint8) (d DynamicFlag) {
	d.IsVisible = flag&1 == 1
	d.VisibilityDidChange = flag&2 == 2
	d.OpacityDidChange = flag&4 == 4
	d.DrawOrderDidChange = flag&8 == 8
	d.RenderOrderDidChange = flag&16 == 16
	d.VertexPositionsDidChange = flag&32 == 32
	d.BlendColorDidChange = flag&64 == 64
	return
}
