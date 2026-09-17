package core

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/aethiopicuschan/cubism-go/internal/core/drawable"
	"github.com/aethiopicuschan/cubism-go/internal/core/moc"
	"github.com/aethiopicuschan/cubism-go/internal/core/parameter"
	"github.com/aethiopicuschan/cubism-go/internal/strings"
	"github.com/aethiopicuschan/cubism-go/internal/utils"
	"github.com/ebitengine/purego"
)

type Core struct {
	lib                           uintptr
	csmGetVersion                 func() uint32
	csmReviveMocInPlace           func(uintptr, uint) uintptr
	csmGetSizeofModel             func(uintptr) uint
	csmInitializeModelInPlace     func(uintptr, uintptr, uint) uintptr
	csmUpdateModel                func(uintptr)
	csmReadCanvasInfo             func(uintptr, uintptr, uintptr, uintptr)
	csmGetParameterCount          func(uintptr) int
	csmGetParameterIds            func(uintptr) **byte
	csmGetParameterTypes          func(uintptr) *int32
	csmGetParameterMinimumValues  func(uintptr) *float32
	csmGetParameterMaximumValues  func(uintptr) *float32
	csmGetParameterDefaultValues  func(uintptr) *float32
	csmGetParameterValues         func(uintptr) *float32
	csmGetPartCount               func(uintptr) int
	csmGetPartIds                 func(uintptr) **byte
	csmGetPartOpacities           func(uintptr) *float32
	csmGetDrawableCount           func(uintptr) int
	csmGetDrawableIds             func(uintptr) **byte
	csmGetDrawableConstantFlags   func(uintptr) *uint8
	csmGetDrawableDynamicFlags    func(uintptr) *uint8
	csmGetDrawableTextureIndices  func(uintptr) *int32
	csmGetDrawableRenderOrders    func(uintptr) *int32
	csmGetDrawableOpacities       func(uintptr) *float32
	csmGetDrawableMaskCounts      func(uintptr) *int32
	csmGetDrawableMasks           func(uintptr) **int32
	csmGetDrawableVertexCounts    func(uintptr) *int32
	csmGetDrawableVertexPositions func(uintptr) **drawable.Vector2
	csmGetDrawableVertexUvs       func(uintptr) **drawable.Vector2
	csmGetDrawableIndexCounts     func(uintptr) *int32
	csmGetDrawableIndices         func(uintptr) **uint16
	csmResetDrawableDynamicFlags  func(uintptr)
	csmHasMocConsistency          func(uintptr, uint) int
}

func NewCore(lib uintptr) (c *Core, err error) {
	c = new(Core)
	c.lib = lib
	purego.RegisterLibFunc(&c.csmGetVersion, lib, "csmGetVersion")
	purego.RegisterLibFunc(&c.csmReviveMocInPlace, lib, "csmReviveMocInPlace")
	purego.RegisterLibFunc(&c.csmGetSizeofModel, lib, "csmGetSizeofModel")
	purego.RegisterLibFunc(&c.csmInitializeModelInPlace, lib, "csmInitializeModelInPlace")
	purego.RegisterLibFunc(&c.csmUpdateModel, lib, "csmUpdateModel")
	purego.RegisterLibFunc(&c.csmReadCanvasInfo, lib, "csmReadCanvasInfo")
	purego.RegisterLibFunc(&c.csmGetParameterCount, lib, "csmGetParameterCount")
	purego.RegisterLibFunc(&c.csmGetParameterIds, lib, "csmGetParameterIds")
	purego.RegisterLibFunc(&c.csmGetParameterTypes, lib, "csmGetParameterTypes")
	purego.RegisterLibFunc(&c.csmGetParameterMinimumValues, lib, "csmGetParameterMinimumValues")
	purego.RegisterLibFunc(&c.csmGetParameterMaximumValues, lib, "csmGetParameterMaximumValues")
	purego.RegisterLibFunc(&c.csmGetParameterDefaultValues, lib, "csmGetParameterDefaultValues")
	purego.RegisterLibFunc(&c.csmGetParameterValues, lib, "csmGetParameterValues")
	purego.RegisterLibFunc(&c.csmGetPartCount, lib, "csmGetPartCount")
	purego.RegisterLibFunc(&c.csmGetPartIds, lib, "csmGetPartIds")
	purego.RegisterLibFunc(&c.csmGetPartOpacities, lib, "csmGetPartOpacities")
	purego.RegisterLibFunc(&c.csmGetDrawableCount, lib, "csmGetDrawableCount")
	purego.RegisterLibFunc(&c.csmGetDrawableIds, lib, "csmGetDrawableIds")
	purego.RegisterLibFunc(&c.csmGetDrawableConstantFlags, lib, "csmGetDrawableConstantFlags")
	purego.RegisterLibFunc(&c.csmGetDrawableDynamicFlags, lib, "csmGetDrawableDynamicFlags")
	purego.RegisterLibFunc(&c.csmGetDrawableTextureIndices, lib, "csmGetDrawableTextureIndices")
	purego.RegisterLibFunc(&c.csmGetDrawableRenderOrders, lib, "csmGetDrawableRenderOrders")
	purego.RegisterLibFunc(&c.csmGetDrawableOpacities, lib, "csmGetDrawableOpacities")
	purego.RegisterLibFunc(&c.csmGetDrawableMaskCounts, lib, "csmGetDrawableMaskCounts")
	purego.RegisterLibFunc(&c.csmGetDrawableMasks, lib, "csmGetDrawableMasks")
	purego.RegisterLibFunc(&c.csmGetDrawableVertexCounts, lib, "csmGetDrawableVertexCounts")
	purego.RegisterLibFunc(&c.csmGetDrawableVertexPositions, lib, "csmGetDrawableVertexPositions")
	purego.RegisterLibFunc(&c.csmGetDrawableVertexUvs, lib, "csmGetDrawableVertexUvs")
	purego.RegisterLibFunc(&c.csmGetDrawableIndexCounts, lib, "csmGetDrawableIndexCounts")
	purego.RegisterLibFunc(&c.csmGetDrawableIndices, lib, "csmGetDrawableIndices")
	purego.RegisterLibFunc(&c.csmResetDrawableDynamicFlags, lib, "csmResetDrawableDynamicFlags")
	purego.RegisterLibFunc(&c.csmHasMocConsistency, lib, "csmHasMocConsistency")
	return
}

// Load moc3 and return moc.Moc
func (c *Core) LoadMoc(path string) (moc moc.Moc, err error) {
	// Read the moc3
	moc.MocBuffer, err = os.ReadFile(path)
	if err != nil {
		return
	}
	// Check the consistency
	consistency := c.csmHasMocConsistency(uintptr(unsafe.Pointer(&moc.MocBuffer[0])), uint(len(moc.MocBuffer)))
	if consistency != 1 {
		err = fmt.Errorf("moc3 is not consistent")
		return
	}
	// Load the moc3
	moc.MocPtr = c.csmReviveMocInPlace(uintptr(unsafe.Pointer(&moc.MocBuffer[0])), uint(len(moc.MocBuffer)))
	if moc.MocPtr == 0 {
		err = fmt.Errorf("failed to revive moc3")
		return
	}
	// Get size
	size := c.csmGetSizeofModel(moc.MocPtr)
	if size == 0 {
		err = fmt.Errorf("failed to get size of model")
		return
	}
	// Initialize the model
	moc.ModelBuffer = make([]byte, size)
	moc.ModelPtr = c.csmInitializeModelInPlace(moc.MocPtr, uintptr(unsafe.Pointer(&moc.ModelBuffer[0])), size)
	if moc.ModelPtr == 0 {
		err = fmt.Errorf("failed to initialize model")
		return
	}

	return
}

// Get version
func (c *Core) GetVersion() string {
	raw := c.csmGetVersion()
	return utils.ParseVersion(raw)
}

// Get dynamic flags
func (c *Core) GetDynamicFlags(modelPtr uintptr) (rs []drawable.DynamicFlag) {
	count := c.csmGetDrawableCount(modelPtr)
	raw := unsafe.Slice(c.csmGetDrawableDynamicFlags(modelPtr), count)
	for _, flag := range raw {
		rs = append(rs, drawable.ParseDynamicFlag(flag))
	}
	return
}

// Get opacities
func (c *Core) GetOpacities(modelPtr uintptr) (rs []float32) {
	count := c.csmGetDrawableCount(modelPtr)
	rs = unsafe.Slice(c.csmGetDrawableOpacities(modelPtr), count)
	return
}

// Get vertex positions
func (c *Core) GetVertexPositions(modelPtr uintptr) (vps [][]drawable.Vector2) {
	count := c.csmGetDrawableCount(modelPtr)
	// 頂点の数
	vertexCounts := unsafe.Slice(c.csmGetDrawableVertexCounts(modelPtr), count)
	posPtrs := unsafe.Slice(c.csmGetDrawableVertexPositions(modelPtr), count)
	for i := range count {
		vertexCount := vertexCounts[i]
		positions := unsafe.Slice(posPtrs[i], int(vertexCount))
		vps = append(vps, positions)
	}
	return
}

// Get Drawables
// Since all the information is gathered, the cost is high. It is expected to be called only once initially
func (c *Core) GetDrawables(modelPtr uintptr) (ds []drawable.Drawable) {
	count := c.csmGetDrawableCount(modelPtr)

	constantFlags := make([]drawable.ConstantFlag, 0)
	raw := unsafe.Slice(c.csmGetDrawableConstantFlags(modelPtr), count)
	for _, flag := range raw {
		constantFlags = append(constantFlags, drawable.ParseConstantFlag(flag))
	}

	dynamicFlags := c.GetDynamicFlags(modelPtr)

	textureIndices := unsafe.Slice(c.csmGetDrawableTextureIndices(modelPtr), count)

	opacities := c.GetOpacities(modelPtr)

	vertexCounts := unsafe.Slice(c.csmGetDrawableVertexCounts(modelPtr), count)

	vertexPositions := make([][]drawable.Vector2, 0)
	vertexUvs := make([][]drawable.Vector2, 0)
	posPtrs := unsafe.Slice(c.csmGetDrawableVertexPositions(modelPtr), count)
	uvPtrs := unsafe.Slice(c.csmGetDrawableVertexUvs(modelPtr), count)
	for i := range count {
		vertexCount := vertexCounts[i]
		positions := unsafe.Slice(posPtrs[i], int(vertexCount))
		vertexPositions = append(vertexPositions, positions)
		uvs := unsafe.Slice(uvPtrs[i], int(vertexCount))
		vertexUvs = append(vertexUvs, uvs)
	}

	// Size of the array of corresponding numbers for the polygon
	indexCounts := unsafe.Slice(c.csmGetDrawableIndexCounts(modelPtr), count)
	// Array of corresponding numbers for the polygon
	indices := make([][]uint16, 0)
	indicesPtrs := unsafe.Slice(c.csmGetDrawableIndices(modelPtr), count)
	for i := range count {
		indexCount := indexCounts[i]
		indices = append(indices, unsafe.Slice(indicesPtrs[i], int(indexCount)))
	}

	// Number of masks
	maskCounts := unsafe.Slice(c.csmGetDrawableMaskCounts(modelPtr), count)
	// Masks
	masks := make([][]int32, 0)
	maskPtrs := unsafe.Slice(c.csmGetDrawableMasks(modelPtr), count)
	for i := range count {
		maskCount := maskCounts[i]
		masks = append(masks, unsafe.Slice(maskPtrs[i], int(maskCount)))
	}

	// ID
	idPtrs := unsafe.Slice(c.csmGetDrawableIds(modelPtr), count)
	ids := make([]string, 0)
	for i := range count {
		ids = append(ids, strings.GoString(idPtrs[i]))
	}

	// Pack into a structure
	for i := range count {
		d := drawable.Drawable{
			Id:              ids[i],
			Texture:         textureIndices[i],
			VertexPositions: vertexPositions[i],
			VertexUvs:       vertexUvs[i],
			VertexIndices:   indices[i],
			ConstantFlag:    constantFlags[i],
			DynamicFlag:     dynamicFlags[i],
			Opacity:         opacities[i],
			Masks:           masks[i],
		}
		ds = append(ds, d)
	}
	return
}

// Get parameters
func (c *Core) GetParameters(modelPtr uintptr) (parameters []parameter.Parameter) {
	count := c.csmGetParameterCount(modelPtr)
	idPtrs := unsafe.Slice(c.csmGetParameterIds(modelPtr), count)
	minPtr := c.csmGetParameterMinimumValues(modelPtr)
	mins := unsafe.Slice(minPtr, count)
	maxPtr := c.csmGetParameterMaximumValues(modelPtr)
	maxs := unsafe.Slice(maxPtr, count)
	defPtr := c.csmGetParameterDefaultValues(modelPtr)
	defs := unsafe.Slice(defPtr, count)
	valPtr := c.csmGetParameterValues(modelPtr)
	vals := unsafe.Slice(valPtr, count)
	for i := range count {
		parameter := parameter.Parameter{
			Id:      strings.GoString(idPtrs[i]),
			Minimum: mins[i],
			Maximum: maxs[i],
			Default: defs[i],
			Current: vals[i],
		}
		parameters = append(parameters, parameter)
	}
	return
}

// Get parameter value
func (c *Core) GetParameterValue(modelPtr uintptr, id string) float32 {
	count := c.csmGetParameterCount(modelPtr)
	idPtrs := unsafe.Slice(c.csmGetParameterIds(modelPtr), count)
	valPtr := c.csmGetParameterValues(modelPtr)
	vals := unsafe.Slice(valPtr, count)
	for i := range count {
		_id := strings.GoString(idPtrs[i])
		if _id == id {
			return vals[i]
		}
	}
	return 0
}

// Set parameter value
func (c *Core) SetParameterValue(modelPtr uintptr, id string, value float32) {
	count := c.csmGetParameterCount(modelPtr)
	idPtrs := unsafe.Slice(c.csmGetParameterIds(modelPtr), count)
	vals := unsafe.Slice(c.csmGetParameterValues(modelPtr), count)
	for i := range count {
		if strings.GoString(idPtrs[i]) == id {
			vals[i] = value
			return
		}
	}
}

// Get the part IDs
func (c *Core) GetPartIds(modelPtr uintptr) (ids []string) {
	count := c.csmGetPartCount(modelPtr)
	idPtrs := unsafe.Slice(c.csmGetPartIds(modelPtr), count)
	for i := range count {
		ids = append(ids, strings.GoString(idPtrs[i]))
	}
	return
}

// Set the part's opacity
func (c *Core) SetPartOpacity(modelPtr uintptr, id string, value float32) {
	ids := c.GetPartIds(modelPtr)
	opacities := unsafe.Slice(c.csmGetPartOpacities(modelPtr), len(ids))
	for i, _id := range ids {
		if _id == id {
			opacities[i] = value
			return
		}
	}
}

// Get the drawing order
// The index of the n-th drawable to be drawn can be obtained with rs[n].
func (c *Core) GetSortedDrawableIndices(modelPtr uintptr) (rs []int) {
	// Drawableの数
	count := c.csmGetDrawableCount(modelPtr)
	// 描画順を取得する
	ptr := c.csmGetDrawableRenderOrders(modelPtr)
	rawIndices := unsafe.Slice(ptr, count)
	rs = make([]int, count)
	for i, order := range rawIndices {
		rs[order] = i
	}
	return
}

// Get the canvas info
func (c *Core) GetCanvasInfo(modelPtr uintptr) (size drawable.Vector2, origin drawable.Vector2, pixelsPerUnit float32) {
	c.csmReadCanvasInfo(modelPtr, uintptr(unsafe.Pointer(&size)), uintptr(unsafe.Pointer(&origin)), uintptr(unsafe.Pointer(&pixelsPerUnit)))
	return
}

// Update the model
func (c *Core) Update(modelPtr uintptr) {
	c.csmResetDrawableDynamicFlags(modelPtr)
	c.csmUpdateModel(modelPtr)
}
