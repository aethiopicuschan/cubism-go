package utils_test

import (
	"testing"

	"github.com/aethiopicuschan/cubism-go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	testcases := []struct {
		name   string
		src    uint32
		expect string
	}{
		{
			name:   "1.0.0",
			src:    16777216,
			expect: "1.0.0",
		},
		{
			name:   "2.3.4",
			src:    33751044,
			expect: "2.3.4",
		},
		{
			name:   "5.0.0",
			src:    83886080,
			expect: "5.0.0",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			got := utils.ParseVersion(testcase.src)
			assert.Equal(t, testcase.expect, got)
		})
	}
}
