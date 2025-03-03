// exiftool_windows.go
//go:build windows

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed assets/windows/ExifTool.exe
var exiftoolWindowsFS embed.FS

func listExiftoolFiles() {
	exiftoolFilesFS, err := fs.Sub(exiftoolWindowsFS, "exiftool_files")
	if err != nil {
		fmt.Printf("Failed to get exiftool_files sub FS: %v\n", err)
		return
	}

	err = fs.WalkDir(exiftoolFilesFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fmt.Println(path)
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to list exiftool_files: %v\n", err)
	}
}

func extractPlatformSpecificExiftool(tempDir string) (string, error) {
	exiftoolName := "exiftool_windows.exe"
	exiftoolPath := filepath.Join(tempDir, exiftoolName)
	exiftoolData, err := exiftoolWindowsFS.ReadFile(exiftoolName)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded exiftool: %v", err)
	}

	err = os.WriteFile(exiftoolPath, exiftoolData, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to write exiftool to disk: %v", err)
	}

	return exiftoolPath, nil
}
