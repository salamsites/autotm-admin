package models

type Slider struct {
	ID          int64
	ImagePathTM []string
	ImagePathEN []string
	ImagePathRU []string
	UploadIdTM  string
	UploadIdEN  string
	UploadIdRU  string
	Platform    string
}
