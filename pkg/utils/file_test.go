package utils_test

import (
	"testing"

	"github.com/dpx/dpx/pkg/mock/assert"
	"github.com/dpx/dpx/pkg/utils"
)

func Test_IsSubPath(t *testing.T) {
	tests := []struct {
		path1    string
		path2    string
		expected bool
	}{
		{"/usr/local/bin", "/usr/local", true},
		{"/usr/local", "/usr/local/bin", true},
		{"/usr/local", "/usr/local", true},
		{"/usr/local", "/tmp", false},
	}

	for _, test := range tests {
		r := utils.IsSubPath(test.path1, test.path2)

		assert.Equal(t, test.expected, r)
	}
}
