//go:build darwin

package main

import (
	_ "embed"
)

const windowsFileName = "windows.txt"

//go:embed assets/windows.txt
var windowsFile []byte
