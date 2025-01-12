package renderer

import (
	_ "embed"
	"image"
	"image/color"
	_ "image/png"

	"github.com/aethiopicuschan/cubism-go"
	"github.com/aethiopicuschan/cubism-go/renderer/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed  mask.kage
var maskShaderSrc []byte

type Renderer struct {
	fb, mb, surface *ebiten.Image
	textureMap      map[string]*ebiten.Image
	model           *cubism.Model
	drawables       []cubism.Drawable
	vertices        [][]ebiten.Vertex
	maskShader      *ebiten.Shader
	final           image.Rectangle
}

// Constructor for the [Renderer] struct
func NewRenderer(model *cubism.Model) (r *Renderer, err error) {
	modelPtr := model.GetMoc().ModelPtr
	core := model.GetCore()
	size, _, _ := core.GetCanvasInfo(modelPtr)
	m := make(map[string]*ebiten.Image)
	ts := model.GetTextures()
	for _, t := range ts {
		img, _, err := ebitenutil.NewImageFromFile(t)
		if err != nil {
			return nil, err
		}
		m[t] = img
	}
	shader, err := ebiten.NewShader(maskShaderSrc)
	if err != nil {
		return
	}
	r = &Renderer{
		fb:         ebiten.NewImage(int(size.X), int(size.Y)),
		mb:         ebiten.NewImage(int(size.X), int(size.Y)),
		surface:    ebiten.NewImage(int(size.X), int(size.Y)),
		textureMap: m,
		model:      model,
		maskShader: shader,
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
	for _, index := range sortedIndices {
		d := r.drawables[index]
		if !d.DynamicFlag.IsVisible {
			continue
		}
		vertices := r.vertices[index]
		if len(d.Masks) > 0 {
			r.fb.Fill(color.RGBA{0, 0, 0, 0})
			r.mb.Fill(color.RGBA{0, 0, 0, 0})
			changed := false
			for _, maskIndex := range d.Masks {
				mask := r.drawables[maskIndex]
				if !mask.DynamicFlag.VertexPositionsDidChange {
					continue
				}
				maskOptions := &colorm.DrawTrianglesOptions{}
				maskColorM := colorm.ColorM{}
				maskColorM.Scale(0, 0, 0, 1)
				maskOptions.AntiAlias = true
				colorm.DrawTriangles(r.mb, r.vertices[maskIndex], mask.VertexIndices, r.textureMap[mask.Texture], maskColorM, maskOptions)
				changed = true
			}
			if changed {
				r.fb.DrawTriangles(vertices, d.VertexIndices, r.textureMap[d.Texture], &ebiten.DrawTrianglesOptions{})
				options := &ebiten.DrawRectShaderOptions{}
				options.Images[0] = r.mb
				options.Images[1] = r.fb
				r.surface.DrawRectShader(r.fb.Bounds().Dx(), r.fb.Bounds().Dy(), r.maskShader, options)
			}
		} else {
			colorM := colorm.ColorM{}
			colorM.Scale(1, 1, 1, float64(d.Opacity))
			options := &colorm.DrawTrianglesOptions{}
			options.AntiAlias = true
			colorm.DrawTriangles(r.surface, vertices, d.VertexIndices, r.textureMap[d.Texture], colorM, options)
		}
	}

	// Draw
	screen.DrawImage(r.surface, last_options)
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
