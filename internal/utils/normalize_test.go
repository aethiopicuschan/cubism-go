package utils_test

import (
	"fmt"
	"testing"

	"github.com/aethiopicuschan/cubism-go/internal/utils"
)

func assertFloat32(exp float32, got float32) (err error) {
	if exp != got {
		err = fmt.Errorf("Expected %f, got %f", exp, got)
	}
	return
}

func TestNormalize(t *testing.T) {
	testcases := []struct {
		x      float32
		n      float32
		m      float32
		expect float32
	}{
		{
			x:      5,
			n:      0,
			m:      10,
			expect: 0,
		},
		{
			x:      0,
			n:      0,
			m:      10,
			expect: -1,
		},
		{
			x:      10,
			n:      0,
			m:      10,
			expect: 1,
		},
	}

	for _, testcase := range testcases {
		got := utils.Normalize(testcase.x, testcase.n, testcase.m)
		if err := assertFloat32(testcase.expect, got); err != nil {
			t.Error(err)
		}
	}
}
