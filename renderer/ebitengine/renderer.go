package renderer

import (
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/aethiopicuschan/cubism-go"
	"github.com/aethiopicuschan/cubism-go/renderer/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Ebitengine buffers are premultiplied: mask coverage must scale both RGB
// and alpha, otherwise clipped pixels still contribute color when composited.
var (
	normalMaskBlend = ebiten.Blend{
		BlendFactorSourceRGB:        ebiten.BlendFactorZero,
		BlendFactorDestinationRGB:   ebiten.BlendFactorSourceAlpha,
		BlendOperationRGB:           ebiten.BlendOperationAdd,
		BlendFactorSourceAlpha:      ebiten.BlendFactorZero,
		BlendFactorDestinationAlpha: ebiten.BlendFactorSourceAlpha,
		BlendOperationAlpha:         ebiten.BlendOperationAdd,
	}
	invertedMaskBlend = ebiten.Blend{
		BlendFactorSourceRGB:        ebiten.BlendFactorZero,
		BlendFactorDestinationRGB:   ebiten.BlendFactorOneMinusSourceAlpha,
		BlendOperationRGB:           ebiten.BlendOperationAdd,
		BlendFactorSourceAlpha:      ebiten.BlendFactorZero,
		BlendFactorDestinationAlpha: ebiten.BlendFactorOneMinusSourceAlpha,
		BlendOperationAlpha:         ebiten.BlendOperationAdd,
	}
)

type Renderer struct {
	fb, mb, surface *ebiten.Image
	textureMap      map[string]*ebiten.Image
	model           *cubism.Model
	drawables       []cubism.Drawable
	vertices        [][]ebiten.Vertex
	final           image.Rectangle
}

// Options for constructing a [Renderer]
type RendererOption struct {
	maxResolution int
}

// Cap the largest dimension of the internal offscreen buffers (fb/mb/surface)
// to at most max pixels, downscaling proportionally if the model's native
// canvas size (as reported by the Cubism core) is larger. This has no effect
// on visual correctness (the model is still rendered at the requested output
// size via the final scale in [Renderer.Draw]), but some models report a
// canvas size far larger than any reasonable display resolution (e.g. Mao's
// 5800x8400), which would otherwise force every offscreen fill/blend/clear
// operation to run at that native, needlessly huge resolution every frame.
// Pass 0 to disable capping and always use the model's native canvas size.
func WithMaxResolution(max int) func(*RendererOption) {
	return func(o *RendererOption) {
		o.maxResolution = max
	}
}

// Constructor for the [Renderer] struct
func NewRenderer(model *cubism.Model, opts ...func(*RendererOption)) (r *Renderer, err error) {
	opt := &RendererOption{
		maxResolution: 2048,
	}
	for _, o := range opts {
		o(opt)
	}
	modelPtr := model.GetMoc().ModelPtr
	core := model.GetCore()
	size, _, _ := core.GetCanvasInfo(modelPtr)
	width, height := size.X, size.Y
	if opt.maxResolution > 0 {
		if largest := float32(math.Max(float64(width), float64(height))); largest > float32(opt.maxResolution) {
			scale := float32(opt.maxResolution) / largest
			width *= scale
			height *= scale
		}
	}
	m := make(map[string]*ebiten.Image)
	ts := model.GetTextures()
	for _, t := range ts {
		img, _, err := ebitenutil.NewImageFromFile(t)
		if err != nil {
			return nil, err
		}
		m[t] = img
	}
	r = &Renderer{
		fb:         ebiten.NewImage(int(width), int(height)),
		mb:         ebiten.NewImage(int(width), int(height)),
		surface:    ebiten.NewImage(int(width), int(height)),
		textureMap: m,
		model:      model,
	}
	return
}

// Update the renderer
func (r *Renderer) Update() error {
	r.model.Update(1.0 / float64(ebiten.TPS()))
	r.drawables = r.model.GetDrawables()
	vertices := make([][]ebiten.Vertex, 0)
	for _, d := range r.drawables {
		v := make([]ebiten.Vertex, 0)
		for i := 0; i < len(d.VertexPositions); i++ {
			v = append(v, ebiten.Vertex{
				DstX:   (d.VertexPositions[i].X + 1) * float32(r.surface.Bounds().Dx()) / 2,
				DstY:   (d.VertexPositions[i].Y*-1 + 1) * float32(r.surface.Bounds().Dy()) / 2,
				SrcX:   d.VertexUvs[i].X * float32(r.textureMap[d.Texture].Bounds().Dx()),
				SrcY:   (1 - d.VertexUvs[i].Y) * float32(r.textureMap[d.Texture].Bounds().Dy()),
				ColorR: 1,
				ColorG: 1,
				ColorB: 1,
				ColorA: 1,
			})
		}
		vertices = append(vertices, v)
	}
	r.vertices = vertices
	return nil
}

// Options for drawing
type DrawOption struct {
	hidden     bool
	scale      float64
	x, y       float64
	background color.Color
}

// Prevent rendering to the final screen
func WithHidden() func(*DrawOption) {
	return func(o *DrawOption) {
		o.hidden = true
	}
}

// Set the scale
func WithScale(scale float64) func(*DrawOption) {
	return func(o *DrawOption) {
		o.scale = scale
	}
}

// Set the position
func WithPosition(x, y float64) func(*DrawOption) {
	return func(o *DrawOption) {
		o.x = x
		o.y = y
	}
}

// Set the background color
func WithBackground(c color.Color) func(*DrawOption) {
	return func(o *DrawOption) {
		o.background = c
	}
}

// Draw the renderer
func (r *Renderer) Draw(screen *ebiten.Image, opts ...func(*DrawOption)) {
	opt := &DrawOption{
		hidden:     false,
		scale:      1,
		x:          0,
		y:          0,
		background: color.Transparent,
	}
	for _, o := range opts {
		o(opt)
	}

	last_options := &ebiten.DrawImageOptions{}
	// First, adjust to the screen size
	screenWidth, screenHeight := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	surfaceWidth, surfaceHeight := float64(r.surface.Bounds().Dx()), float64(r.surface.Bounds().Dy())
	last_options.GeoM.Scale(screenHeight/screenWidth, 1)
	last_options.GeoM.Scale(screenWidth/surfaceWidth, screenHeight/surfaceHeight)
	// Apply the scale options
	last_options.GeoM.Scale(opt.scale, opt.scale)
	// Align the horizontal axis to the center
	width := screenWidth * (screenHeight / screenWidth) * opt.scale
	height := screenHeight * opt.scale
	x := screenWidth/2 - width/2 + opt.x
	y := screenHeight/2 - height/2 + opt.y
	last_options.GeoM.Translate(x, y)
	r.final = image.Rect(int(x), int(y), int(x+width), int(y+height))
	// Set Alpha
	last_options.ColorScale.SetA(r.model.GetOpacity())

	if opt.hidden {
		return
	}

	r.surface.Fill(opt.background)
	sortedIndices := r.model.GetSortedIndices()

	// Consecutive non-masked, visible drawables that share the same texture
	// are batched into a single DrawTriangles call. Each individual draw call
	// carries substantial fixed overhead (command encoding, driver bridging),
	// so submitting ~225 tiny draw calls per frame (one per part) is far more
	// expensive than a handful of larger batched ones. Masked drawables can't
	// join a batch since they require their own clip/composite step, so a
	// pending batch is flushed whenever one is encountered.
	var batchVertices []ebiten.Vertex
	var batchIndices []uint16
	var batchTexture *ebiten.Image
	flushBatch := func() {
		if len(batchIndices) == 0 {
			return
		}
		r.surface.DrawTriangles(batchVertices, batchIndices, batchTexture, &ebiten.DrawTrianglesOptions{})
		batchVertices = batchVertices[:0]
		batchIndices = batchIndices[:0]
		batchTexture = nil
	}

	for _, index := range sortedIndices {
		d := r.drawables[index]
		if !d.DynamicFlag.IsVisible {
			continue
		}
		vertices := r.vertices[index]
		if len(d.Masks) > 0 {
			flushBatch()
			// Only the area actually touched by this drawable and its masks needs
			// to be cleared/composited. Operating on the full canvas here is very
			// expensive for models with large canvases (e.g. Mao's 5800x8400),
			// since it turns every masked part into a full-buffer clear + shader
			// pass regardless of how small the part actually is on screen.
			bounds := verticesBounds(vertices)
			for _, maskIndex := range d.Masks {
				bounds = bounds.Union(verticesBounds(r.vertices[maskIndex]))
			}
			bounds = bounds.Intersect(r.fb.Bounds())
			if bounds.Empty() {
				continue
			}
			subFb := r.fb.SubImage(bounds).(*ebiten.Image)
			subMb := r.mb.SubImage(bounds).(*ebiten.Image)
			subFb.Clear()
			subMb.Clear()
			// Every mask assigned to this drawable must be redrawn every
			// frame it is visible, regardless of whether that particular
			// mask's own vertex positions happen to have changed since the
			// last frame: mb is fully cleared above, so skipping a mask here
			// would leave its contribution missing from the union entirely
			// (not just "stale"). This previously caused e.g. eyelids (whose
			// own clip mask is mostly static) to intermittently fail to
			// render while blinking, letting the eyes poke through them.
			for _, maskIndex := range d.Masks {
				mask := r.drawables[maskIndex]
				maskOptions := &ebiten.DrawTrianglesOptions{}
				r.mb.DrawTriangles(r.vertices[maskIndex], mask.VertexIndices, r.textureMap[mask.Texture], maskOptions)
			}
			r.fb.DrawTriangles(opacityVertices(vertices, d.Opacity), d.VertexIndices, r.textureMap[d.Texture], &ebiten.DrawTrianglesOptions{})
			r.compositeMask(bounds, d.ConstantFlag.IsInvertedMask)
		} else {
			texture := r.textureMap[d.Texture]
			if batchTexture != texture {
				flushBatch()
				batchTexture = texture
			}
			// uint16 indices can only address up to 65536 vertices per draw
			// call; flush and start a new batch before that would overflow.
			if len(batchVertices)+len(vertices) > 65536 {
				flushBatch()
				batchTexture = texture
			}
			base := uint16(len(batchVertices))
			batchVertices = append(batchVertices, opacityVertices(vertices, d.Opacity)...)
			for _, idx := range d.VertexIndices {
				batchIndices = append(batchIndices, base+idx)
			}
		}
	}
	flushBatch()

	// Draw
	screen.DrawImage(r.surface, last_options)
}

func (r *Renderer) compositeMask(bounds image.Rectangle, inverted bool) {
	subFb := r.fb.SubImage(bounds).(*ebiten.Image)
	subMb := r.mb.SubImage(bounds).(*ebiten.Image)
	options := &ebiten.DrawImageOptions{Blend: normalMaskBlend}
	if inverted {
		options.Blend = invertedMaskBlend
	}
	// DrawImage rebases its source to (0, 0), but destination subimages
	// retain their original coordinates.
	options.GeoM.Translate(float64(bounds.Min.X), float64(bounds.Min.Y))
	subFb.DrawImage(subMb, options)
	blitOptions := &ebiten.DrawImageOptions{}
	blitOptions.GeoM.Translate(float64(bounds.Min.X), float64(bounds.Min.Y))
	r.surface.DrawImage(subFb, blitOptions)
}

// Apply an opacity value to a copy of vertices' alpha channel. Vertices with
// opacity 1 (the common case) are returned unmodified to avoid needless
// allocation/copies.
func opacityVertices(vertices []ebiten.Vertex, opacity float32) []ebiten.Vertex {
	if opacity == 1 {
		return vertices
	}
	out := make([]ebiten.Vertex, len(vertices))
	for i, v := range vertices {
		v.ColorA *= opacity
		out[i] = v
	}
	return out
}

// Compute the destination-space bounding rectangle of a set of vertices,
// used to limit mask buffer clears/composites to the area actually in use.
func verticesBounds(vertices []ebiten.Vertex) image.Rectangle {
	if len(vertices) == 0 {
		return image.Rectangle{}
	}
	minX, minY := vertices[0].DstX, vertices[0].DstY
	maxX, maxY := vertices[0].DstX, vertices[0].DstY
	for _, v := range vertices[1:] {
		if v.DstX < minX {
			minX = v.DstX
		}
		if v.DstX > maxX {
			maxX = v.DstX
		}
		if v.DstY < minY {
			minY = v.DstY
		}
		if v.DstY > maxY {
			maxY = v.DstY
		}
	}
	return image.Rect(int(math.Floor(float64(minX))), int(math.Floor(float64(minY))), int(math.Ceil(float64(maxX)))+1, int(math.Ceil(float64(maxY)))+1)
}

// Get the model set in the renderer
func (r *Renderer) GetModel() *cubism.Model {
	return r.model
}

// Perform collision detection
func (r *Renderer) IsHit(x, y int, id string) (hit bool, err error) {
	// Out of bounds
	if r.final.Min.X > x || x > r.final.Max.X || r.final.Min.Y > y || y > r.final.Max.Y {
		return
	}

	// Get the Drawable
	d, err := r.model.GetDrawable(id)
	if err != nil {
		return
	}

	// Rectangular range
	var left, right, top, bottom float32
	left = float32(r.surface.Bounds().Dx())
	top = float32(r.surface.Bounds().Dy())

	// Get the rectangle representing the range of the Drawable
	for i := 0; i < len(d.VertexPositions); i++ {
		v := d.VertexPositions[i]
		if v.X < left {
			left = v.X
		}
		if v.X > right {
			right = v.X
		}
		if v.Y < top {
			top = v.Y
		}
		if v.Y > bottom {
			bottom = v.Y
		}
	}

	// Convert to local coordinates
	localX := utils.Normalize(float32(x), float32(r.final.Min.X), float32(r.final.Max.X))
	localY := utils.Normalize(float32(y), float32(r.final.Min.Y), float32(r.final.Max.Y)) * -1

	if left <= localX && localX <= right && top <= localY && localY <= bottom {
		hit = true
	}

	return
}
