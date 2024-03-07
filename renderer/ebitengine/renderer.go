package renderer

import (
	_ "embed"
	"image"
	"image/color"

	"github.com/aethiopicuschan/cubism-go"
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

// 描画オプション
type DrawOption struct {
	hidden bool
	scale  float64
	x, y   float64
}

// 最終的なscreenに対する描画を行わないようにする
func WithHidden() func(*DrawOption) {
	return func(o *DrawOption) {
		o.hidden = true
	}
}

// スケールを設定する
func WithScale(scale float64) func(*DrawOption) {
	return func(o *DrawOption) {
		o.scale = scale
	}
}

// 位置を設定する
func WithPosition(x, y float64) func(*DrawOption) {
	return func(o *DrawOption) {
		o.x = x
		o.y = y
	}
}

// 描画する
func (r *Renderer) Draw(screen *ebiten.Image, opts ...func(*DrawOption)) {
	opt := &DrawOption{
		hidden: false,
		scale:  1,
		x:      0,
		y:      0,
	}
	for _, o := range opts {
		o(opt)
	}
	screen.Fill(color.White)
	r.surface.Fill(color.RGBA{0, 255, 0, 255})
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
			for _, maskIndex := range d.Masks {
				mask := r.drawables[maskIndex]
				maskOptions := &colorm.DrawTrianglesOptions{}
				maskColorM := colorm.ColorM{}
				maskColorM.Scale(0, 0, 0, 1)
				maskOptions.AntiAlias = true
				colorm.DrawTriangles(r.mb, r.vertices[maskIndex], mask.VertexIndices, r.textureMap[mask.Texture], maskColorM, maskOptions)
			}
			r.fb.DrawTriangles(vertices, d.VertexIndices, r.textureMap[d.Texture], &ebiten.DrawTrianglesOptions{})
			options := &ebiten.DrawRectShaderOptions{}
			options.Images[0] = r.mb
			options.Images[1] = r.fb
			r.surface.DrawRectShader(r.fb.Bounds().Dx(), r.fb.Bounds().Dy(), r.maskShader, options)
		} else {
			colorM := colorm.ColorM{}
			colorM.Scale(1, 1, 1, float64(d.Opacity))
			options := &colorm.DrawTrianglesOptions{}
			options.AntiAlias = true
			colorm.DrawTriangles(r.surface, vertices, d.VertexIndices, r.textureMap[d.Texture], colorM, options)
		}
	}
	options := &ebiten.DrawImageOptions{}
	// まず、画面サイズに合わせる
	screenWidth, screenHeight := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	surfaceWidth, surfaceHeight := float64(r.surface.Bounds().Dx()), float64(r.surface.Bounds().Dy())
	options.GeoM.Scale(screenHeight/screenWidth, 1)
	options.GeoM.Scale(screenWidth/surfaceWidth, screenHeight/surfaceHeight)
	// スケールオプションを適用
	options.GeoM.Scale(opt.scale, opt.scale)
	// 横軸を中央に合わせる
	width := screenWidth * (screenHeight / screenWidth) * opt.scale
	height := screenHeight * opt.scale
	x := screenWidth/2 - width/2
	y := screenHeight/2 - height/2
	options.GeoM.Translate(x+opt.x, y+opt.y)

	// アルファ値
	options.ColorScale.SetA(r.model.Opacity)
	// 描画
	if !opt.hidden {
		screen.DrawImage(r.surface, options)
	}
}

func (r *Renderer) GetModel() *cubism.Model {
	return r.model
}
