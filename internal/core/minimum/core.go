package minimum

import (
	"github.com/aethiopicuschan/cubism-go/internal/utils"
	"github.com/ebitengine/purego"
)

type Core struct {
	csmGetVersion func() uint32
}

func NewCore(lib uintptr) (c Core, err error) {
	purego.RegisterLibFunc(&c.csmGetVersion, lib, "csmGetVersion")
	return
}

func (c Core) GetVersion() string {
	raw := c.csmGetVersion()
	return utils.ParseVersion(raw)
}
