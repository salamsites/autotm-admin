package dtos

type V1BrandDTO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	LogoPath string `json:"logo_path"`
}

type V1BrandModelDTO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	LogoPath string `json:"logo_path"`
	BrandID  int64  `json:"brand_id"`
}
