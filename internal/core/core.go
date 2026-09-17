package core

import (
	"fmt"
	"strings"

	core_5 "github.com/aethiopicuschan/cubism-go/internal/core/core_5"
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
	// Cubism 5 API is stable across the whole 5.x.x range, so this check unnecessarily rejects all newer 5.x.x versions that do work fine with the existing bindings.
	if strings.HasPrefix(version, "5.") {
		c, err = core_5.NewCore(l)
		return
	}
	err = fmt.Errorf("unsupported version: %s", version)
	return
}
