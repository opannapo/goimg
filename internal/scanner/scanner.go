package scanner

import (
	"fmt"
	"goimg/internal/models"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/barasher/go-exiftool"
)

var imageExtensions = []string{
	".jpg", ".jpeg", ".png", ".tiff", ".tif", ".webp", ".heic", ".heif", ".raw", ".cr2", ".nef",
	".mov", ".mp4", ".m4v", // video
}

func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}

func ScanImages(dir string) ([]models.Metadata, error) {
	var results []models.Metadata

	et, err := exiftool.NewExiftool()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize exiftool: %w", err)
	}
	defer et.Close()

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// ignore subdir
			if path != dir {
				return filepath.SkipDir
			}
			return nil
		}
		if !IsImageFile(path) {
			return nil
		}

		meta, err := extractMetadata(et, path)
		if err != nil {
			fmt.Printf("⚠️  Skip %s: %v\n", path, err)
			return nil
		}

		// --- MOVING START ---
		userComment, _ := meta.AllFields["UserComment"].(string)
		deviceMaker, _ := meta.AllFields["DeviceManufacturer"].(string)
		handler, _ := meta.AllFields["HandlerDescription"].(string) // Specifically for video, Apple frequently utilizes 'HandlerDescription' or 'Encoder' tags

		makeLower := strings.ToLower(meta.CameraMake)
		makerLower := strings.ToLower(deviceMaker)
		commentLower := strings.ToLower(userComment)
		handlerLower := strings.ToLower(handler)
		fileExt := strings.ToLower(filepath.Ext(path))

		// Detect whether this is an Apple product (Photo or Video)
		isApple := strings.Contains(makeLower, "apple") ||
			strings.Contains(makerLower, "apple") ||
			strings.Contains(handlerLower, "apple")

		if isApple {
			destFolder := "IPHONE CAMERA" // Default

			// Cek Screenshot (PNG) or Screen Recording (Video)
			isScreenshot := commentLower == "screenshot" || fileExt == ".png"

			if isScreenshot {
				destFolder = "IPHONE SCREENSHOT"
			} else if fileExt == ".mov" || fileExt == ".mp4" {
				destFolder = "IPHONE VIDEO"
			}

			newPath, err := moveFile(path, destFolder)
			if err != nil {
				fmt.Printf("❌ Error moving %s: %v\n", path, err)
			} else {
				meta.FilePath = newPath
				fmt.Printf("🎥 [%s] Moved: %s\n", destFolder, filepath.Base(newPath))
			}
		}
		// --- MOVING END ---

		results = append(results, meta)
		return nil
	})

	return results, err
}

func moveFile(sourcePath, destDir string) (string, error) {
	// Create folder if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}

	fileName := filepath.Base(sourcePath)
	destPath := filepath.Join(destDir, fileName)

	//Relocate file (os.Rename safely preserves metadata)
	err := os.Rename(sourcePath, destPath)
	if err != nil {
		return "", err
	}

	return destPath, nil
}

func extractMetadata(et *exiftool.Exiftool, filePath string) (models.Metadata, error) {
	fileInfos := et.ExtractMetadata(filePath)
	if len(fileInfos) == 0 {
		return models.Metadata{FilePath: filePath}, nil
	}

	fileInfo := fileInfos[0]
	if fileInfo.Err != nil {
		return models.Metadata{FilePath: filePath}, fileInfo.Err
	}

	var size int64
	if val, ok := fileInfo.Fields["FileSize"].(int64); ok {
		size = val
	} else if val, ok := fileInfo.Fields["FileSize"].(float64); ok {
		size = int64(val)
	}

	meta := models.Metadata{
		FilePath:  filePath,
		FileSize:  size,
		FileType:  getField(&fileInfo, "FileType"),
		AllFields: fileInfo.Fields,
	}

	date := getField(&fileInfo, "DateTimeOriginal")
	if date == "" {
		date = getField(&fileInfo, "CreateDate")
	}
	if date == "" {
		date = getField(&fileInfo, "FileModifyDate")
	}

	meta.CreateDate = date
	meta.Dimensions = getField(&fileInfo, "ImageSize")
	meta.CameraMake = getField(&fileInfo, "Make")
	meta.CameraModel = getField(&fileInfo, "Model")
	meta.ISO = getField(&fileInfo, "ISO")
	meta.Aperture = getField(&fileInfo, "Aperture")
	meta.Shutter = getField(&fileInfo, "ExposureTime")
	meta.FocalLen = getField(&fileInfo, "FocalLength")

	if getField(&fileInfo, "GPSLatitude") != "" {
		meta.GPS = map[string]any{
			"latitude":  getField(&fileInfo, "GPSLatitude"),
			"longitude": getField(&fileInfo, "GPSLongitude"),
			"altitude":  getField(&fileInfo, "GPSAltitude"),
		}
	}

	return meta, nil
}

func getField(fileInfo *exiftool.FileMetadata, fieldName string) string {
	if val, exists := fileInfo.Fields[fieldName]; exists {
		return fmt.Sprintf("%v", val)
	}
	return ""
}
