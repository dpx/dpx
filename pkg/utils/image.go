package utils

import (
	"path"
	"regexp"
	"strings"
)

var regexImageName = regexp.MustCompile(`(?i)[^a-z_-]`)

func NormalizeImageName(image string) string {
	s := strings.Split(image, ":")
	if len(s) > 1 {
		image = s[0]
	}

	return regexImageName.ReplaceAllString(image, "")
}

func GetContainerName(image, dir string) string {
	// use dir as container name
	if dir == "" {
		dir = GetCwd()
	}

	name := path.Base(dir)

	// just in case where we're at `/` rootname
	if name == "" || strings.HasPrefix(name, ".") {
		name = "dpx"
	} else if name == "dpx.yml" {
		name = path.Dir(dir)
	}

	// add suffix
	name = strings.Join([]string{NormalizeImageName(name), NormalizeImageName(image)}, "-")

	return name
}
