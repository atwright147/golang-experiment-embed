// exiftool_darwin.go
//go:build darwin || linux

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed assets/nix/*
var exiftoolDarwinFS embed.FS

func listExiftoolFiles() {
	exiftoolFilesFS, err := fs.Sub(exiftoolDarwinFS, "assets/exiftool_files")
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
	exiftoolName := "exiftool"
	exiftoolFilesDir := "assets/nix"

	exiftoolPath := filepath.Join(tempDir, exiftoolFilesDir, exiftoolName)

	exiftoolFilesFS, err := fs.Sub(exiftoolDarwinFS, exiftoolFilesDir)
	if err != nil {
		return "", fmt.Errorf("failed to get exiftool_files sub FS: %v", err)
	}

	err = fs.WalkDir(exiftoolFilesFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		destPath := filepath.Join(tempDir, exiftoolFilesDir, path)
		destDir := filepath.Dir(destPath)
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			err = os.MkdirAll(destDir, 0755)
			if err != nil {
				return err
			}
		}

		fileData, err := fs.ReadFile(exiftoolFilesFS, path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, fileData, 0644)
	})
	if err != nil {
		return "", fmt.Errorf("failed to extract exiftool_files: %v", err)
	}

	err = os.Chmod(exiftoolPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to make exiftool executable: %v", err)
	}

	return exiftoolPath, nil
}
