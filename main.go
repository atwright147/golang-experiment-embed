package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	// listExiftoolFiles()

	exiftoolPath, err := extractExiftool()
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

func extractExiftool() (string, error) {
	return extractPlatformSpecificExiftool() // Calls platform specific file.
}
