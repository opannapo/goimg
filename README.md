# GoImg-Organizer

--- 
### Description
**GoImg-Organizer** is a high-performance CLI tool built with Go designed to automatically scan, extract metadata, and organize Apple device media (iPhone/iPad). It leverages `ExifTool` to identify files and intelligently categorize them into specific folders based on their origin and type.

### Key Features
*   **Smart Categorization**: Automatically sorts files into:
    *   `IPHONE CAMERA`: Original photos captured via camera.
    *   `IPHONE SCREENSHOT`: Static screen captures (PNG).
    *   `IPHONE VIDEO`: Camera recordings and screen recordings (MOV/MP4).
*   **Metadata Extraction**: Deep-scans files using `ExifTool` to detect Apple-specific tags (e.g., `HandlerDescription` or `Encoder`).
*   **Non-Recursive Safety**: Only processes files in the target directory to prevent accidental modification of subfolders.
*   **JSON Reporting**: Generates a detailed audit log of all processed metadata in a timestamped JSON file.
*   **Metadata Preservation**: Moves files using `os.Rename` to ensure file integrity and original timestamps are kept intact.

### Installation
1. Ensure you have [Go](https://golang.org/doc/install) and [ExifTool](https://exiftool.org/) installed.
2. Usage
   ```bash
   ./bin/goimg scan <target-directory>
