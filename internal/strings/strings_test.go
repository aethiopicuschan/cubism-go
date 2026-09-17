package strings_test

import (
	"testing"

	"github.com/aethiopicuschan/cubism-go/internal/strings"
	"github.com/stretchr/testify/require"
)

func TestGoString(t *testing.T) {
	require.Empty(t, strings.GoString(nil))
	for _, text := range []string{"", "ParamEyeLOpen", "first\x00ignored"} {
		buffer := append([]byte(text), 0)
		want := text
		if text == "first\x00ignored" {
			want = "first"
		}
		got := strings.GoString(&buffer[0])
		require.Equal(t, want, got)
		buffer[0] = 'X'
		require.Equal(t, want, got)
	}
}
