package dtos

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
