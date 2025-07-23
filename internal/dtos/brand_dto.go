package dtos

type CreateBodyTypeReq struct {
	NameTM    string `json:"name_tm"`
	NameEN    string `json:"name_en"`
	NameRU    string `json:"name_ru"`
	ImagePath string `json:"image_path"`
	Category  string `json:"category"`
}

type UpdateBodyTypeReq struct {
	ID        int64  `json:"id"`
	NameTM    string `json:"name_tm"`
	NameEN    string `json:"name_en"`
	NameRU    string `json:"name_ru"`
	ImagePath string `json:"image_path"`
	Category  string `json:"category"`
}

type BodyType struct {
	ID        int64  `json:"id"`
	NameTM    string `json:"name_tm"`
	NameEN    string `json:"name_en"`
	NameRU    string `json:"name_ru"`
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

type CreateModelReq struct {
	Name       string `json:"name"`
	BrandID    int64  `json:"brand_id"`
	BodyTypeID int64  `json:"body_type_id"`
}

type UpdateModelReq struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	BrandID    int64  `json:"brand_id"`
	BodyTypeID int64  `json:"body_type_id"`
}

type Model struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	LogoPath     string `json:"logo_path"`
	BrandID      int64  `json:"brand_id"`
	BrandName    string `json:"brand_name"`
	BodyTypeID   int64  `json:"body_type_id"`
	BodyTypeName string `json:"body_type_name"`
}

type ModelResult struct {
	Models []Model `json:"models"`
	Count  int64   `json:"count"`
}
