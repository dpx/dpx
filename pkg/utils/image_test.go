package utils_test

import (
	"testing"

	"github.com/dpx/dpx/pkg/mock/assert"
	"github.com/dpx/dpx/pkg/utils"
)

func Test_GetContainerName(t *testing.T) {
	tests := []struct {
		image string
		dir   string
		name  string
	}{
		{"go", "", "utils-go"},
		{"go", ".", "dpx-go"},
		{"go", "..", "dpx-go"},
		{"go", "app", "app-go"},
		{"go", "app/dpx.yml", "app-go"},
		{"go", "app/test.json", "testjson-go"},
	}

	for _, test := range tests {
		r := utils.GetContainerName(test.image, test.dir)

		assert.Equal(t, test.name, r)
	}
}
