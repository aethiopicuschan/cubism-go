package renderer

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func CompositeMask(fb, mb, surface *ebiten.Image, bounds image.Rectangle, inverted bool) {
	r := &Renderer{fb: fb, mb: mb, surface: surface}
	r.compositeMask(bounds, inverted)
}
