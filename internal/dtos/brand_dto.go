package dtos

type CreateBodyTypeReq struct {
	Name      string `json:"name" binding:"required"`
	ImagePath string `json:"image_path"`
}

type UpdateBodyTypeReq struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ImagePath string `json:"image_path"`
}

type BodyType struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ImagePath string `json:"image_path"`
	Category  string `json:"category"`
}

type BodyTypeResult struct {
	BodyTypes []BodyType `json:"body_types"`
	Count     int64      `json:"count"`
}
type CreateBrandReq struct {
	Name       string   `json:"name"`
	LogoPath   string   `json:"logo_path"`
	Categories []string `json:"categories" binding:"required,dive,oneof=auto moto truck"`
}

type UpdateBrandReq struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	LogoPath   string   `json:"logo_path"`
	Categories []string `json:"categories" binding:"required,dive,oneof=auto moto truck"`
}
type Brand struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	LogoPath   string   `json:"logo_path"`
	Categories []string `json:"categories"`
}

type BrandResult struct {
	Brands []Brand `json:"brands"`
	Count  int64   `json:"count"`
}

type CreateBrandModelReq struct {
	Name    string `json:"name"`
	BrandID int64  `json:"brand_id"`
}

type UpdateBrandModelReq struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	BrandID int64  `json:"brand_id"`
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
