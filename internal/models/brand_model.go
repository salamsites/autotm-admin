package models

type Brand struct {
	ID       int64
	Name     string
	LogoPath string
}

type BrandResult struct {
	Brands []Brand
	Count  int64
}

type ID struct {
	ID int64 `json:"id"`
}

type BrandModel struct {
	ID        int64
	Name      string
	LogoPath  string
	BrandID   int64
	BrandName string
}

type BrandModelResult struct {
	BrandModels []BrandModel
	Count       int64
}
