package dtos

type UploadImage struct {
	UploadID string   `json:"upload_id"`
	Sizes    []string `json:"sizes"`
}
