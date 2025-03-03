// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// listExiftoolFiles()

	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}
	tempDir = filepath.Join(tempDir, "Desktop", "temp")

	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	exiftoolPath, err := extractExiftool(tempDir)
	if err != nil {
		log.Fatalf("Failed to extract exiftool: %v", err)
	}

	imagePath := "image.ARW" // Replace with your image path

	cmd := exec.Command(exiftoolPath, "-Make", "-Model", imagePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute exiftool: %v, output: %s", err, output)
	}

	fmt.Println(string(output))
}

func extractExiftool(tempDir string) (string, error) {
	return extractPlatformSpecificExiftool(tempDir) // Calls platform specific file.
}
