package renderer_test

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/aethiopicuschan/cubism-go/renderer/ebitengine"
	"github.com/hajimehoshi/ebiten/v2"
)

type testGame struct {
	run func()
}

func (g *testGame) Update() error {
	g.run()
	return ebiten.Termination
}

func (*testGame) Draw(*ebiten.Image) {}

func (*testGame) Layout(int, int) (int, int) { return 32, 32 }

func TestMain(m *testing.M) {
	code := 1
	if err := ebiten.RunGame(&testGame{run: func() { code = m.Run() }}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(code)
}

func TestCompositeMask(t *testing.T) {
	for _, bounds := range []image.Rectangle{
		image.Rect(0, 0, 4, 4),
		image.Rect(12, 16, 16, 20),
	} {
		for _, inverted := range []bool{false, true} {
			t.Run(fmt.Sprintf("bounds=%v/inverted=%v", bounds, inverted), func(t *testing.T) {
				fb := ebiten.NewImage(32, 32)
				mb := ebiten.NewImage(32, 32)
				surface := ebiten.NewImage(32, 32)
				defer fb.Deallocate()
				defer mb.Deallocate()
				defer surface.Deallocate()
				// Reuse buffers across open, half-closed, closed and reopened frames.
				for _, alpha := range []uint8{255, 128, 0, 255} {
					fb.Clear()
					mb.Clear()
					surface.Fill(color.RGBA{0, 0, 255, 255})
					fb.SubImage(bounds).(*ebiten.Image).Fill(color.RGBA{128, 0, 0, 128})
					// Black mask texels ensure coverage comes from alpha, not RGB.
					mb.SubImage(bounds).(*ebiten.Image).Fill(color.RGBA{0, 0, 0, alpha})
					renderer.CompositeMask(fb, mb, surface, bounds, inverted)
					coverage := int(alpha)
					if inverted {
						coverage = 255 - coverage
					}
					red := (128*coverage + 127) / 255
					want := color.RGBA{uint8(red), 0, uint8(255 - red), 255}
					for y := range 32 {
						for x := range 32 {
							expected := color.RGBA{0, 0, 255, 255}
							if image.Pt(x, y).In(bounds) {
								expected = want
							}
							got := surface.At(x, y).(color.RGBA)
							for c, pair := range [][2]uint8{{got.R, expected.R}, {got.G, expected.G}, {got.B, expected.B}, {got.A, expected.A}} {
								diff := int(pair[0]) - int(pair[1])
								if diff < -1 || diff > 1 {
									t.Fatalf("mask alpha=%d pixel=(%d,%d) channel=%d: got %v, want %v", alpha, x, y, c, got, expected)
								}
							}
						}
					}
				}
			})
		}
	}
}
