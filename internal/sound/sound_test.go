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
		src    string
		expect string
	}{
		{
			src:    "test.wav",
			expect: "wav",
		},
		{
			src:    "test.wave",
			expect: "wav",
		},
		{
			src:    "test.mp3",
			expect: "mp3",
		},
	}

	for _, testcase := range testcases {
		got, err := sound.DetectFormat(testcase.src)
		if err != nil {
			t.Error(err)
		}
		if err := assert(testcase.expect, got); err != nil {
			t.Error(err)
		}
	}
}
