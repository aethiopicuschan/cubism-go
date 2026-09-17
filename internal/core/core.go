package core

import (
	"fmt"
	"strings"

	core_5 "github.com/aethiopicuschan/cubism-go/internal/core/core_5"
	core_6 "github.com/aethiopicuschan/cubism-go/internal/core/core_6"
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
	switch {
	case strings.HasPrefix(version, "5."):
		c, err = core_5.NewCore(l)
		return
	case strings.HasPrefix(version, "6."):
		c, err = core_6.NewCore(l)
		return
	}
	err = fmt.Errorf("unsupported version: %s", version)
	return
}
