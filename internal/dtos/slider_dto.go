package dtos

type CreateSliderReq struct {
	ImagePathTM []string `json:"image_path_tm"`
	ImagePathEN []string `json:"image_path_en"`
	ImagePathRU []string `json:"image_path_ru"`
	UploadIdTM  string   `json:"upload_id_tm"`
	UploadIdEN  string   `json:"upload_id_en"`
	UploadIdRU  string   `json:"upload_id_ru"`
	Platform    string   `json:"platform"`
}

type UpdateSliderReq struct {
	ID          int64    `json:"id"`
	ImagePathTM []string `json:"image_path_tm"`
	ImagePathEN []string `json:"image_path_en"`
	ImagePathRU []string `json:"image_path_ru"`
	UploadIdTM  string   `json:"upload_id_tm"`
	UploadIdEN  string   `json:"upload_id_en"`
	UploadIdRU  string   `json:"upload_id_ru"`
	Platform    string   `json:"platform"`
}
type Slider struct {
	ID          int64    `json:"id"`
	ImagePathTM []string `json:"image_path_tm"`
	ImagePathEN []string `json:"image_path_en"`
	ImagePathRU []string `json:"image_path_ru"`
	UploadIdTM  string   `json:"upload_id_tm"`
	UploadIdEN  string   `json:"upload_id_en"`
	UploadIdRU  string   `json:"upload_id_ru"`
	Platform    string   `json:"platform"`
}

type SliderResult struct {
	Sliders []Slider `json:"sliders"`
	Count   int64    `json:"count"`
}
