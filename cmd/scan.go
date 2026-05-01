package cmd

import (
	"encoding/json"
	"fmt"
	"goimg/internal/models"
	"goimg/internal/scanner"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func NewScanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan [directory]",
		Short: "Scan all image files and extract metadata.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 && args[0] != "" {
				dir = args[0]
			}

			absPath, _ := filepath.Abs(dir)
			fmt.Printf("\n📂 Target Directory: %s\n", absPath)
			fmt.Println(strings.Repeat("-", 40))

			images, err := scanner.ScanImages(absPath)
			if err != nil {
				return fmt.Errorf("error performing scan %w", err)
			}

			printResults(images)
			reportResults(images)
			return nil
		},
	}
	return cmd
}

func printResults(images []models.Metadata) {
	fmt.Printf("📸 Found %d image files\n\n", len(images))

	for i, img := range images {
		fmt.Printf("=== FILE %d: %s ===\n", i+1, img.FilePath)
		fmt.Printf("📁 Size: %s\n", formatBytes(img.FileSize))
		fmt.Printf("📋 Type: %s\n", img.FileType)
		fmt.Printf("📐 Dimensions: %s\n", img.Dimensions)
		fmt.Printf("📅 Created: %s\n", img.CreateDate)
		fmt.Printf("📷 Camera: %s %s\n", img.CameraMake, img.CameraModel)
		fmt.Printf("⚙️  ISO: %s | F/%s | %s | %smm\n",
			img.ISO, img.Aperture, img.Shutter, img.FocalLen)

		if len(img.GPS) > 0 {
			fmt.Printf("📍 GPS: Lat:%s, Lng:%s, Alt:%s\n",
				img.GPS["latitude"], img.GPS["longitude"], img.GPS["altitude"])
		}

		fmt.Println("--- RAW METADATA (JSON) ---")
		jsonData, _ := json.MarshalIndent(img.AllFields, "", "  ")
		fmt.Println(string(jsonData))
		fmt.Println("\n" + strings.Repeat("=", 80) + "\n")
	}
}

func reportResults(images []models.Metadata) {
	jsonData, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v\n", err)
		return
	}

	fileName := time.Now().Format("20060102-150405") + ".json"

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		fmt.Printf("Failed to create JSON: %v\n", err)
		return
	}

	fmt.Printf("\n✨ Full report has been saved to: %s\n", fileName)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
