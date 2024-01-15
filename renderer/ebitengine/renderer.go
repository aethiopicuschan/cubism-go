package renderer

import (
	_ "embed"
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
	model           cubism.Model
	drawables       []cubism.Drawable
	vertices        [][]ebiten.Vertex
	maskShader      *ebiten.Shader
}

func NewRenderer(model cubism.Model) (r *Renderer, err error) {
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

func (r *Renderer) Draw(screen *ebiten.Image) {
	width := screen.Bounds().Dx()
	height := screen.Bounds().Dy()
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
	// 中央に表示する
	options.GeoM.Scale(float64(height)/float64(width), 1)
	options.GeoM.Scale(float64(screen.Bounds().Dx())/float64(r.surface.Bounds().Dx()), float64(screen.Bounds().Dy())/float64(r.surface.Bounds().Dy()))
	finalWidth := float64(screen.Bounds().Dx()) * (float64(height) / float64(width))
	options.GeoM.Translate(float64(screen.Bounds().Dx())/2-finalWidth/2, 0)
	options.ColorScale.SetA(r.model.Opacity)
	screen.DrawImage(r.surface, options)
}
