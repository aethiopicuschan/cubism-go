package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/aethiopicuschan/cubism-go"
	renderer "github.com/aethiopicuschan/cubism-go/renderer/ebitengine"
	"github.com/aethiopicuschan/cubism-go/sound/normal"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	Name   = "Haru"
	Width  = 2880
	Height = 1800
)

type Game struct {
	ow, oh   int
	tapId    int
	renderer *renderer.Renderer
}

func (g *Game) Update() (err error) {
	g.renderer.Update()
	x, y := ebiten.CursorPosition()
	if x < 0 || y < 0 || x > g.ow || y > g.oh {
		return
	}
	if !ebiten.IsFocused() {
		return
	}
	hitareas := g.renderer.GetModel().GetHitAreas()
	hitted := false
	for _, hitarea := range hitareas {
		hit, err := g.renderer.IsHit(x, y, hitarea.Id)
		if err != nil {
			return err
		}
		if hit {
			hitted = true
		}
	}
	if hitted {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.renderer.GetModel().StopMotion(g.tapId)
			g.tapId = g.renderer.GetModel().PlayMotion("TapBody", 0, false)
		}
	} else if ebiten.CursorShape() == ebiten.CursorShapePointer {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
	return
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	g.renderer.Draw(screen, renderer.WithBackground(color.RGBA{0, 255, 0, 255}))
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
	// Set function for playing sound
	csm.LoadSound = normal.LoadSound
	model, err := csm.LoadModel(fmt.Sprintf("Resources/%s/%s.model3.json", Name, Name))
	if err != nil {
		log.Fatal(err)
	}
	// Play idle motion
	model.PlayMotion("Idle", 0, true)
	renderer, err := renderer.NewRenderer(model)
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
