#include <stdint.h>

static int32_t render_orders[] = {1, 0};
static int32_t zero_ints[] = {0, 0};
static float values[] = {0.25f, 0.75f};
static uint8_t flags[] = {1, 1};
static const char *ids[] = {"first", "second"};
static const void *empty_arrays[] = {0, 0};

uint32_t csmGetVersion(void) { return (CORE_MAJOR << 24) | (1 << 16) | 2; }
void *csmReviveMocInPlace(void *address, uint32_t size) { return address; }
uint32_t csmGetSizeofModel(const void *moc) { return 64; }
void *csmInitializeModelInPlace(const void *moc, void *address, uint32_t size) {
    return address;
}
int32_t csmHasMocConsistency(void *address, uint32_t size) { return size > 0; }
void csmUpdateModel(void *model) {}
void csmResetDrawableDynamicFlags(void *model) {}
void csmReadCanvasInfo(const void *model, float *size, float *origin, float *ppu) {
    size[0] = 100;
    size[1] = 200;
    origin[0] = 50;
    origin[1] = 100;
    *ppu = 10;
}

#if CORE_MAJOR == 5
const int32_t *csmGetDrawableRenderOrders(const void *model) { return render_orders; }
#else
const int32_t *csmGetRenderOrders(const void *model) { return render_orders; }
int32_t csmGetOffscreenCount(const void *model) { return OFFSCREEN_COUNT; }
#endif

int32_t csmGetParameterCount(const void *model) { return 2; }
const char **csmGetParameterIds(const void *model) { return ids; }
const int32_t *csmGetParameterTypes(const void *model) { return zero_ints; }
const float *csmGetParameterMinimumValues(const void *model) { return values; }
const float *csmGetParameterMaximumValues(const void *model) { return values; }
const float *csmGetParameterDefaultValues(const void *model) { return values; }
float *csmGetParameterValues(void *model) { return values; }
int32_t csmGetPartCount(const void *model) { return 2; }
const char **csmGetPartIds(const void *model) { return ids; }
float *csmGetPartOpacities(void *model) { return values; }
int32_t csmGetDrawableCount(const void *model) { return 2; }
const char **csmGetDrawableIds(const void *model) { return ids; }
const uint8_t *csmGetDrawableConstantFlags(const void *model) { return flags; }
const uint8_t *csmGetDrawableDynamicFlags(const void *model) { return flags; }
const int32_t *csmGetDrawableTextureIndices(const void *model) { return zero_ints; }
const float *csmGetDrawableOpacities(const void *model) { return values; }
const int32_t *csmGetDrawableMaskCounts(const void *model) { return zero_ints; }
const int32_t **csmGetDrawableMasks(const void *model) { return (const int32_t **)empty_arrays; }
const int32_t *csmGetDrawableVertexCounts(const void *model) { return zero_ints; }
const void **csmGetDrawableVertexPositions(const void *model) { return empty_arrays; }
const void **csmGetDrawableVertexUvs(const void *model) { return empty_arrays; }
const int32_t *csmGetDrawableIndexCounts(const void *model) { return zero_ints; }
const uint16_t **csmGetDrawableIndices(const void *model) { return (const uint16_t **)empty_arrays; }
