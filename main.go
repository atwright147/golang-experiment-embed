package main

import (
	"embed"
	"fmt"
	"os"
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

	file1, err := files.ReadFile("assets/file1.txt")
	if err != nil {
		fmt.Printf("Failed to read embedded file: %v\n", err)
		return
	}

	err = os.WriteFile(fmt.Sprintf("%s/file1.txt", path), file1, 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
		return
	}

	fmt.Println("File written successfully")
}
