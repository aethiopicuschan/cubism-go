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
	csmGetParameterIds            func(uintptr) uintptr
	csmGetParameterTypes          func(uintptr) uintptr
	csmGetParameterMinimumValues  func(uintptr) uintptr
	csmGetParameterMaximumValues  func(uintptr) uintptr
	csmGetParameterDefaultValues  func(uintptr) uintptr
	csmGetParameterValues         func(uintptr) uintptr
	csmGetPartCount               func(uintptr) int
	csmGetPartIds                 func(uintptr) uintptr
	csmGetPartOpacities           func(uintptr) uintptr
	csmGetDrawableCount           func(uintptr) int
	csmGetDrawableIds             func(uintptr) uintptr
	csmGetDrawableConstantFlags   func(uintptr) uintptr
	csmGetDrawableDynamicFlags    func(uintptr) uintptr
	csmGetDrawableTextureIndices  func(uintptr) uintptr
	csmGetDrawableRenderOrders    func(uintptr) uintptr
	csmGetDrawableOpacities       func(uintptr) uintptr
	csmGetDrawableMaskCounts      func(uintptr) uintptr
	csmGetDrawableMasks           func(uintptr) uintptr
	csmGetDrawableVertexCounts    func(uintptr) uintptr
	csmGetDrawableVertexPositions func(uintptr) uintptr
	csmGetDrawableVertexUvs       func(uintptr) uintptr
	csmGetDrawableIndexCounts     func(uintptr) uintptr
	csmGetDrawableIndices         func(uintptr) uintptr
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
	raw := unsafe.Slice((*uint8)(unsafe.Pointer(c.csmGetDrawableDynamicFlags(modelPtr))), count)
	for _, flag := range raw {
		rs = append(rs, drawable.ParseDynamicFlag(flag))
	}
	return
}

// Get opacities
func (c *Core) GetOpacities(modelPtr uintptr) (rs []float32) {
	count := c.csmGetDrawableCount(modelPtr)
	rs = unsafe.Slice((*float32)(unsafe.Pointer(c.csmGetDrawableOpacities(modelPtr))), count)
	return
}

// Get vertex positions
func (c *Core) GetVertexPositions(modelPtr uintptr) (vps [][]drawable.Vector2) {
	count := c.csmGetDrawableCount(modelPtr)
	// 頂点の数
	vertexCounts := unsafe.Slice((*int32)(unsafe.Pointer(c.csmGetDrawableVertexCounts(modelPtr))), count)
	posPtr := c.csmGetDrawableVertexPositions(modelPtr)
	for i := 0; i < count; i++ {
		vertexCount := vertexCounts[i]
		positions := unsafe.Slice(*(**drawable.Vector2)(unsafe.Pointer(posPtr + uintptr(i)*unsafe.Sizeof(uintptr(0)))), int(vertexCount))
		vps = append(vps, positions)
	}
	return
}

// Get Drawables
// Since all the information is gathered, the cost is high. It is expected to be called only once initially
func (c *Core) GetDrawables(modelPtr uintptr) (ds []drawable.Drawable) {
	count := c.csmGetDrawableCount(modelPtr)

	constantFlags := make([]drawable.ConstantFlag, 0)
	raw := unsafe.Slice((*uint8)(unsafe.Pointer(c.csmGetDrawableConstantFlags(modelPtr))), count)
	for _, flag := range raw {
		constantFlags = append(constantFlags, drawable.ParseConstantFlag(flag))
	}

	dynamicFlags := c.GetDynamicFlags(modelPtr)

	textureIndices := unsafe.Slice((*int32)(unsafe.Pointer(c.csmGetDrawableTextureIndices(modelPtr))), count)

	opacities := c.GetOpacities(modelPtr)

	vertexCounts := unsafe.Slice((*int32)(unsafe.Pointer(c.csmGetDrawableVertexCounts(modelPtr))), count)

	vertexPositions := make([][]drawable.Vector2, 0)
	vertexUvs := make([][]drawable.Vector2, 0)
	posPtr := c.csmGetDrawableVertexPositions(modelPtr)
	uvPtr := c.csmGetDrawableVertexUvs(modelPtr)
	for i := 0; i < count; i++ {
		vertexCount := vertexCounts[i]
		positions := unsafe.Slice(*(**drawable.Vector2)(unsafe.Pointer(posPtr + uintptr(i)*unsafe.Sizeof(uintptr(0)))), int(vertexCount))
		vertexPositions = append(vertexPositions, positions)
		uvs := unsafe.Slice(*(**drawable.Vector2)(unsafe.Pointer(uvPtr + uintptr(i)*unsafe.Sizeof(uintptr(0)))), int(vertexCount))
		vertexUvs = append(vertexUvs, uvs)
	}

	// Size of the array of corresponding numbers for the polygon
	indexCounts := unsafe.Slice((*int32)(unsafe.Pointer(c.csmGetDrawableIndexCounts(modelPtr))), count)
	// Array of corresponding numbers for the polygon
	indices := make([][]uint16, 0)
	indicesPtr := c.csmGetDrawableIndices(modelPtr)
	for i := 0; i < count; i++ {
		indexCount := indexCounts[i]
		indices = append(indices, unsafe.Slice(*(**uint16)(unsafe.Pointer(indicesPtr + uintptr(i)*unsafe.Sizeof(uintptr(0)))), int(indexCount)))
	}

	// Number of masks
	maskCounts := unsafe.Slice((*int32)(unsafe.Pointer(c.csmGetDrawableMaskCounts(modelPtr))), count)
	// Masks
	masks := make([][]int32, 0)
	maskPtr := c.csmGetDrawableMasks(modelPtr)
	for i := 0; i < count; i++ {
		maskCount := maskCounts[i]
		masks = append(masks, unsafe.Slice(*(**int32)(unsafe.Pointer(maskPtr + uintptr(i)*unsafe.Sizeof(uintptr(0)))), int(maskCount)))
	}

	// ID
	idsPtr := c.csmGetDrawableIds(modelPtr)
	ids := make([]string, 0)
	for i := 0; i < count; i++ {
		ptr := *(**byte)(unsafe.Pointer(idsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
		ids = append(ids, strings.GoString(uintptr(unsafe.Pointer(ptr))))
	}

	// Pack into a structure
	for i := 0; i < count; i++ {
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
	idsPtr := c.csmGetParameterIds(modelPtr)
	minPtr := c.csmGetParameterMinimumValues(modelPtr)
	mins := unsafe.Slice((*float32)(unsafe.Pointer(minPtr)), count)
	maxPtr := c.csmGetParameterMaximumValues(modelPtr)
	maxs := unsafe.Slice((*float32)(unsafe.Pointer(maxPtr)), count)
	defPtr := c.csmGetParameterDefaultValues(modelPtr)
	defs := unsafe.Slice((*float32)(unsafe.Pointer(defPtr)), count)
	valPtr := c.csmGetParameterValues(modelPtr)
	vals := unsafe.Slice((*float32)(unsafe.Pointer(valPtr)), count)
	for i := 0; i < count; i++ {
		ptr := *(**byte)(unsafe.Pointer(idsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
		parameter := parameter.Parameter{
			Id:      strings.GoString(uintptr(unsafe.Pointer(ptr))),
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
	idsPtr := c.csmGetParameterIds(modelPtr)
	valPtr := c.csmGetParameterValues(modelPtr)
	vals := unsafe.Slice((*float32)(unsafe.Pointer(valPtr)), count)
	for i := 0; i < count; i++ {
		ptr := *(**byte)(unsafe.Pointer(idsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
		_id := strings.GoString(uintptr(unsafe.Pointer(ptr)))
		if _id == id {
			return vals[i]
		}
	}
	return 0
}

// Set parameter value
func (c *Core) SetParameterValue(modelPtr uintptr, id string, value float32) {
	count := c.csmGetParameterCount(modelPtr)
	idsPtr := c.csmGetParameterIds(modelPtr)
	valPtr := c.csmGetParameterValues(modelPtr)
	for i := 0; i < count; i++ {
		ptr := *(**byte)(unsafe.Pointer(idsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
		if strings.GoString(uintptr(unsafe.Pointer(ptr))) == id {
			*(*float32)(unsafe.Pointer(valPtr + uintptr(i)*unsafe.Sizeof(float32(0)))) = value
			return
		}
	}
}

// Get the part IDs
func (c *Core) GetPartIds(modelPtr uintptr) (ids []string) {
	count := c.csmGetPartCount(modelPtr)
	idsPtr := c.csmGetPartIds(modelPtr)
	for i := 0; i < count; i++ {
		ptr := *(**byte)(unsafe.Pointer(idsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
		ids = append(ids, strings.GoString(uintptr(unsafe.Pointer(ptr))))
	}
	return
}

// Set the part's opacity
func (c *Core) SetPartOpacity(modelPtr uintptr, id string, value float32) {
	ids := c.GetPartIds(modelPtr)
	ptr := c.csmGetPartOpacities(modelPtr)
	for i, _id := range ids {
		if _id == id {
			*(*float32)(unsafe.Pointer(ptr + uintptr(i)*unsafe.Sizeof(float32(0)))) = value
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
	rawIndices := unsafe.Slice((*int32)(unsafe.Pointer(ptr)), count)
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
