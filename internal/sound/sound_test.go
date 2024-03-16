package sound_test

import (
	"fmt"
	"testing"

	"github.com/aethiopicuschan/cubism-go/internal/sound"
)

func assert(exp string, got string) (err error) {
	if exp != got {
		err = fmt.Errorf("Expected %s, got %s", exp, got)
	}
	return
}

func TestDetectFormat(t *testing.T) {
	testcases := []struct {
		src       string
		expect    string
		expectErr bool
	}{
		{
			src:       "test.wav",
			expect:    "wav",
			expectErr: false,
		},
		{
			src:       "test.wave",
			expect:    "wav",
			expectErr: false,
		},
		{
			src:       "test.mp3",
			expect:    "mp3",
			expectErr: false,
		},
		{
			src:       "test.aac",
			expect:    "",
			expectErr: true,
		},
	}

	for _, testcase := range testcases {
		got, err := sound.DetectFormat(testcase.src)
		if testcase.expectErr {
			if err == nil {
				t.Error("Expected error, got nil")
			}
		} else {
			if err != nil {
				t.Error(err)
			}
		}
		if err := assert(testcase.expect, got); err != nil {
			t.Error(err)
		}
	}
}
