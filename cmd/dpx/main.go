package main

import (
	"fmt"
	"os"

	"github.com/dpx/dpx/pkg/core"
)

func main() {
	if err := NewApp().Run(os.Args); err != nil {
		exit(err)
	}
}

func exit(err error) {
	switch e := err.(type) {
	case *core.ContainerExecErr:
		os.Exit(e.Code)
	default:
		fmt.Println(err)
		os.Exit(1)
	}
}
