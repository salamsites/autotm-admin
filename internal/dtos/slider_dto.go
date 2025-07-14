package dtos

type CreateSliderReq struct {
	ImagePath string `json:"image_path" validate:"required"`
	Title     string `json:"title" validate:"required"`
	Platform  string `json:"platform" validate:"required"`
}

type UpdateSliderReq struct {
	ID        int64  `json:"id"`
	ImagePath string `json:"image_path"`
	Title     string `json:"title"`
	Platform  string `json:"platform"`
}
type Slider struct {
	ID        int64  `json:"id"`
	ImagePath string `json:"image_path"`
	Title     string `json:"title"`
	Platform  string `json:"platform"`
}

type SliderResult struct {
	Sliders []Slider `json:"sliders"`
	Count   int64    `json:"count"`
}
