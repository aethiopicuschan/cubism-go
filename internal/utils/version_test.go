package utils_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/aethiopicuschan/cubism-go/internal/utils"
)

func createVersion(major, minor, patch int) uint32 {
	return (uint32(major) << 24) | (uint32(minor) << 16) | uint32(patch)
}

func assert(exp string, got string) (err error) {
	if exp != got {
		err = fmt.Errorf("Expected %s, got %s", exp, got)
	}
	return
}

func TestParseVersion(t *testing.T) {
	log.Println(createVersion(2, 3, 4))
	testcases := []struct {
		src    uint32
		expect string
	}{
		{
			src:    16777216,
			expect: "1.0.0",
		},
		{
			src:    33751044,
			expect: "2.3.4",
		},
		{
			src:    83886080,
			expect: "5.0.0",
		},
	}

	for _, testcase := range testcases {
		got := utils.ParseVersion(testcase.src)
		if err := assert(testcase.expect, got); err != nil {
			t.Error(err)
		}
	}
}
