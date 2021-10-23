//go:build tools
// +build tools

// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md

package tools

import (
	_ "github.com/suntong/html2md"
)
