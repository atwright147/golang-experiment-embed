package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// exiftoolBinary is declared in platform-specific files (exiftool_darwin.go, exiftool_windows.go)
var exiftoolBinary embed.FS

func getExifToolPath() (string, func(), error) {
	tempDir, err := os.MkdirTemp("", "exiftool")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	var exifToolPath string
	var binaryName string

	if runtime.GOOS == "darwin" {
		binaryName = "exiftool_assets/darwin/exiftool"
	} else if runtime.GOOS == "windows" {
		binaryName = "exiftool_assets/windows/exiftool.exe"
	} else {
		return "", cleanup, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	// Extract the binary
	binaryData, err := exiftoolBinary.ReadFile(binaryName)
	if err != nil {
		return "", cleanup, fmt.Errorf("failed to read embedded binary: %w", err)
	}
	exifToolPath = filepath.Join(tempDir, filepath.Base(binaryName))
	err = os.WriteFile(exifToolPath, binaryData, 0755)
	if err != nil {
		return "", cleanup, fmt.Errorf("failed to write ExifTool binary: %w", err)
	}

	// For macOS, extract additional files
	if runtime.GOOS == "darwin" {
		libDir := filepath.Join(tempDir, "lib")
		log.Printf("exiftool lib dir: %s", libDir)
		err = os.Mkdir(libDir, 0755)
		if err != nil {
			return "", cleanup, fmt.Errorf("failed to create lib directory: %w", err)
		}

		err = extractDir(exiftoolBinary, "exiftool_assets/darwin/exiftool_files", libDir)
		if err != nil {
			return "", cleanup, fmt.Errorf("failed to extract lib files: %w", err)
		}
	}

	return exifToolPath, cleanup, nil
}

func extractDir(fs embed.FS, src, dst string) error {
	entries, err := fs.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := extractDir(fs, srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := fs.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func useExifTool(imagePath string) (string, error) {
	exifToolPath, cleanup, err := getExifToolPath()
	if err != nil {
		return "", err
	}
	defer cleanup()

	cmd := exec.Command(exifToolPath,
		"-thumbnailimage",
		"-b",
		"-w", "%d%f_thumb.jpg",
		imagePath)

	if runtime.GOOS == "darwin" {
		libPath := filepath.Join(filepath.Dir(exifToolPath), "lib")
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("PERL5LIB=%s", libPath))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ExifTool execution failed: %w\nOutput: %s", err, output)
	}

	return string(output), nil
}

// App struct for Wails binding
type App struct{}

// ExtractMetadata is the method exposed to the frontend
func (a *App) ExtractMetadata(imagePath string) string {
	metadata, err := useExifTool(imagePath)
	if err != nil {
		return "ExtractMetadata Error: " + err.Error()
	}
	return metadata
}

func _main() {
	// Example usage
	imagePath := "./DSC01451.ARW"
	metadata, err := useExifTool(imagePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Metadata:", metadata)

	// For Wails, you would typically set up your app here
	// err := wails.Run(&options.App{
	//     Title:  "My ExifTool App",
	//     Width:  1024,
	//     Height: 768,
	//     AssetServer: &assetserver.Options{
	//         Assets: assets,
	//     },
	//     BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
	//     OnStartup:        app.startup,
	//     Bind: []interface{}{
	//         &App{},
	//     },
	// })
	// if err != nil {
	//     println("Error:", err.Error())
	// }
}
