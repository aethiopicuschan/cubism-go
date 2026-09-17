package core

import (
	"fmt"

	core_5 "github.com/aethiopicuschan/cubism-go/internal/core/core_5"
	"github.com/aethiopicuschan/cubism-go/internal/core/moc"
	"github.com/ebitengine/purego"
)

// Core uses the common Core 5 bindings with the Core 6 render-order API.
// Offscreen composition is not supported by the current renderer interface.
type Core struct {
	*core_5.Core
	csmGetOffscreenCount func(uintptr) int32
}

func NewCore(lib uintptr) (*Core, error) {
	common, err := core_5.NewCoreWithRenderOrders(lib, "csmGetRenderOrders")
	if err != nil {
		return nil, err
	}
	c := &Core{Core: common}
	purego.RegisterLibFunc(&c.csmGetOffscreenCount, lib, "csmGetOffscreenCount")
	return c, nil
}

func (c *Core) LoadMoc(path string) (moc.Moc, error) {
	m, err := c.Core.LoadMoc(path)
	if err != nil {
		return m, err
	}
	count := c.csmGetOffscreenCount(m.ModelPtr)
	if count < 0 {
		return moc.Moc{}, fmt.Errorf("failed to get Core 6 offscreen count")
	}
	if count > 0 {
		return moc.Moc{}, fmt.Errorf("Core 6 model uses %d offscreens: offscreen composition is not supported", count)
	}
	return m, nil
}
