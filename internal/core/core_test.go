//go:build darwin || freebsd || linux

package core_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/aethiopicuschan/cubism-go/internal/core"
	core_5 "github.com/aethiopicuschan/cubism-go/internal/core/core_5"
	core_6 "github.com/aethiopicuschan/cubism-go/internal/core/core_6"
	"github.com/stretchr/testify/require"
)

var (
	_ core.Core = (*core_5.Core)(nil)
	_ core.Core = (*core_6.Core)(nil)
)

func buildCoreLibrary(t *testing.T, major, offscreens int) string {
	t.Helper()
	cc, err := exec.LookPath("cc")
	if err != nil {
		t.Skip("native binding tests require a C compiler")
	}
	path := filepath.Join(t.TempDir(), "core.so")
	args := []string{"-shared", "-fPIC"}
	if runtime.GOOS == "darwin" {
		args = []string{"-dynamiclib"}
	}
	args = append(args, fmt.Sprintf("-DCORE_MAJOR=%d", major),
		fmt.Sprintf("-DOFFSCREEN_COUNT=%d", offscreens), "testdata/core.c", "-o", path)
	output, err := exec.Command(cc, args...).CombinedOutput()
	require.NoError(t, err, "%s", output)
	return path
}

func TestCoreVersionSelection(t *testing.T) {
	for _, major := range []int{5, 6, 7, 60} {
		t.Run(fmt.Sprintf("major=%d", major), func(t *testing.T) {
			c, err := core.NewCore(buildCoreLibrary(t, major, 0))
			if major != 5 && major != 6 {
				require.ErrorContains(t, err, "unsupported version:")
				require.Nil(t, c)
				return
			}
			require.NoError(t, err)
			if major == 5 {
				require.IsType(t, &core_5.Core{}, c)
			} else {
				require.IsType(t, &core_6.Core{}, c)
			}
			require.Equal(t, fmt.Sprintf("%d.1.2", major), c.GetVersion())
			path := filepath.Join(t.TempDir(), "model.moc3")
			require.NoError(t, os.WriteFile(path, []byte{1}, 0600))
			m, err := c.LoadMoc(path)
			require.NoError(t, err)
			require.NotZero(t, m.ModelPtr)
			require.Equal(t, []int{1, 0}, c.GetSortedDrawableIndices(m.ModelPtr))
			require.Len(t, c.GetDrawables(m.ModelPtr), 2)
			require.Len(t, c.GetVertexPositions(m.ModelPtr), 2)
			require.Equal(t, float32(0.75), c.GetParameterValue(m.ModelPtr, "second"))
			c.SetParameterValue(m.ModelPtr, "second", 0.5)
			require.Equal(t, float32(0.5), c.GetParameters(m.ModelPtr)[1].Current)
			require.Equal(t, []string{"first", "second"}, c.GetPartIds(m.ModelPtr))
			c.SetPartOpacity(m.ModelPtr, "first", 0.5)
			size, origin, ppu := c.GetCanvasInfo(m.ModelPtr)
			require.Equal(t, float32(200), size.Y)
			require.Equal(t, float32(50), origin.X)
			require.Equal(t, float32(10), ppu)
			c.Update(m.ModelPtr)
			runtime.KeepAlive(m)
		})
	}
}

func TestCore6Offscreens(t *testing.T) {
	for _, count := range []int{-1, 1} {
		t.Run(fmt.Sprintf("count=%d", count), func(t *testing.T) {
			c, err := core.NewCore(buildCoreLibrary(t, 6, count))
			require.NoError(t, err)
			path := filepath.Join(t.TempDir(), "model.moc3")
			require.NoError(t, os.WriteFile(path, []byte{1}, 0600))
			m, err := c.LoadMoc(path)
			if count < 0 {
				require.ErrorContains(t, err, "failed to get Core 6 offscreen count")
			} else {
				require.ErrorContains(t, err, "offscreen composition is not supported")
			}
			require.Zero(t, m.ModelPtr)
			_, err = c.LoadMoc(filepath.Join(t.TempDir(), "missing.moc3"))
			require.ErrorIs(t, err, os.ErrNotExist)
		})
	}
}
