package models

type Metadata struct {
	FilePath    string         `json:"file_path"`
	FileSize    int64          `json:"file_size"`
	FileType    string         `json:"file_type"`
	Dimensions  string         `json:"dimensions"`
	CreateDate  string         `json:"create_date"`
	CameraMake  string         `json:"camera_make"`
	CameraModel string         `json:"camera_model"`
	ISO         string         `json:"iso"`
	Aperture    string         `json:"aperture"`
	Shutter     string         `json:"shutter_speed"`
	FocalLen    string         `json:"focal_length"`
	GPS         map[string]any `json:"gps,omitempty"`
	AllFields   map[string]any `json:"all_fields,omitempty"`
}
