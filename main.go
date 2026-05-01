package main

import (
	"fmt"
	"goimg/cmd"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "goimg",
		Short: "A Go-based CLI tool that uses ExifTool to automatically organize Apple photos, screenshots, and videos into categorized folders based on metadata",
	}

	rootCmd.AddCommand(cmd.NewScanCmd())

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
