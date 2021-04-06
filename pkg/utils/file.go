package utils

import (
	"os"
	"strings"
)

// FileExist check if file exists
func FileExist(f string) bool {
	if _, err := os.Stat(f); err != nil {
		return false
	}

	return true
}

// GetCwd returns current working dir
func GetCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("Cannot get working directory")
	}

	return dir
}

func JoinPath(paths ...string) string {
	return strings.Join(paths, string(os.PathListSeparator))
}

// IsTTY term.isTerminal(fd)
func IsTTY() bool {
	fi, _ := os.Stdout.Stat()

	return (fi.Mode() & os.ModeCharDevice) != 0
}

func IsStdIn() bool {
	fi, _ := os.Stdin.Stat()

	return (fi.Mode() & os.ModeCharDevice) == 0
}

// IsSubPath checks if path1 is a subpath of path2 or vise versa.
func IsSubPath(p1, p2 string) bool {
	return strings.Index(p1, p2) > -1 || strings.Index(p2, p1) > -1
}
