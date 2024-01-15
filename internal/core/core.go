package core

import (
	"fmt"

	core_5_0_0 "github.com/aethiopicuschan/cubism-go/internal/core/core_5_0_0"
	"github.com/aethiopicuschan/cubism-go/internal/core/drawable"
	"github.com/aethiopicuschan/cubism-go/internal/core/minimum"
	"github.com/aethiopicuschan/cubism-go/internal/core/moc"
	"github.com/aethiopicuschan/cubism-go/internal/core/parameter"
)

type Core interface {
	LoadMoc(path string) (moc.Moc, error)
	GetVersion() string
	GetDynamicFlags(uintptr) []drawable.DynamicFlag
	GetOpacities(uintptr) []float32
	GetVertexPositions(uintptr) [][]drawable.Vector2
	GetDrawables(uintptr) []drawable.Drawable
	GetParameters(uintptr) []parameter.Parameter
	GetParameterValue(uintptr, string) float32
	SetParameterValue(uintptr, string, float32)
	GetPartIds(uintptr) []string
	SetPartOpacity(uintptr, string, float32)
	GetSortedDrawableIndices(uintptr) []int
	GetCanvasInfo(uintptr) (drawable.Vector2, drawable.Vector2, float32)
	Update(uintptr)
}

func NewCore(lib string) (c Core, err error) {
	l, err := openLibrary(lib)
	if err != nil {
		return
	}
	mc, err := minimum.NewCore(l)
	if err != nil {
		return
	}
	version := mc.GetVersion()
	if version == "5.0.0" {
		c, err = core_5_0_0.NewCore(l)
		return
	}
	err = fmt.Errorf("unsupported version: %s", version)
	return
}
