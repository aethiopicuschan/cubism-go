package utils_test

import (
	"testing"

	"github.com/aethiopicuschan/cubism-go/renderer/utils"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	testcases := []struct {
		name   string
		x      float32
		n      float32
		m      float32
		expect float32
	}{
		{
			name:   "center",
			x:      5,
			n:      0,
			m:      10,
			expect: 0,
		},
		{
			name:   "min",
			x:      0,
			n:      0,
			m:      10,
			expect: -1,
		},
		{
			name:   "max",
			x:      10,
			n:      0,
			m:      10,
			expect: 1,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			got := utils.Normalize(testcase.x, testcase.n, testcase.m)
			assert.Equal(t, testcase.expect, got)
		})
	}
}
