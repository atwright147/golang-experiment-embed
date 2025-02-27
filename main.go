package main

import (
	"embed"
	"fmt"
	"os"
	"runtime"
)

//go:embed assets/*
var files embed.FS

func main() {
	path := "temp"
	err := os.Mkdir(path, 0755)
	if err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return
	}

	// defer os.RemoveAll(path)

	var fileName string

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		fileName = "assets/nix.txt"
	} else if runtime.GOOS == "windows" {
		fileName = "assets/windows.txt"
	} else {
		err = fmt.Errorf("Unsupported OS: %s", runtime.GOOS)
		fmt.Printf("Error: %v\n", err)
		return
	}

	file, err := files.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Failed to read embedded file: %v\n", err)
		return
	}

	err = os.WriteFile(fmt.Sprintf("%s/os_file.txt", path), file, 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
		return
	}

	fmt.Println("File written successfully")
}
