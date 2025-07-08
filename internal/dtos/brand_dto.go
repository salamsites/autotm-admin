package dtos

type Brand struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	LogoPath string `json:"logo_path"`
}

type BrandResult struct {
	Brands []Brand `json:"brands"`
	Count  int64   `json:"count"`
}

type BrandModel struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	LogoPath  string `json:"logo_path"`
	BrandID   int64  `json:"brand_id"`
	BrandName string `json:"brand_name"`
}

type BrandModelResult struct {
	BrandModels []BrandModel `json:"brand_models"`
	Count       int64        `json:"count"`
}
