//go:build darwin || linux

package main

import (
	_ "embed"
)

const nixFileName = "nix.txt"

//go:embed assets/nix.txt
var nixFile []byte
