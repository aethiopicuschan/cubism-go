package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/aethiopicuschan/cubism-go"
	renderer "github.com/aethiopicuschan/cubism-go/renderer/ebitengine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	Name   = "Haru"
	Width  = 2880
	Height = 1800
)

type Game struct {
	ow, oh   int
	renderer *renderer.Renderer
}

func (g *Game) Update() error {
	g.renderer.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	g.renderer.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.ow, g.oh = outsideWidth, outsideHeight
	return outsideWidth, outsideHeight
}

func main() {
	csm, err := cubism.NewCubism("libLive2DCubismCore.dylib")
	if err != nil {
		log.Fatal(err)
	}
	model, err := csm.LoadModel(fmt.Sprintf("Resources/%s/%s.model3.json", Name, Name))
	if err != nil {
		log.Fatal(err)
	}
	renderer, err := renderer.NewRenderer(model)
	model.PlayMotion("Idle", 0)
	if err != nil {
		log.Fatal(err)
	}
	g := &Game{
		renderer: renderer,
	}
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
