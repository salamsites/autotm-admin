package dtos

type CreateSliderReq struct {
	ImagePathTM string `json:"image_path_tm" validate:"required"`
	ImagePathEN string `json:"image_path_en" validate:"required"`
	ImagePathRU string `json:"image_path_ru" validate:"required"`
	Platform    string `json:"platform" validate:"required"`
}

type UpdateSliderReq struct {
	ID          int64  `json:"id"`
	ImagePathTM string `json:"image_path_tm"`
	ImagePathEN string `json:"image_path_en"`
	ImagePathRU string `json:"image_path_ru"`
	Platform    string `json:"platform"`
}
type Slider struct {
	ID          int64  `json:"id"`
	ImagePathTM string `json:"image_path_tm"`
	ImagePathEN string `json:"image_path_en"`
	ImagePathRU string `json:"image_path_ru"`
	Platform    string `json:"platform"`
}

type SliderResult struct {
	Sliders []Slider `json:"sliders"`
	Count   int64    `json:"count"`
}
