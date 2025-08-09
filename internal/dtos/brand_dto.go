package dtos

type ID struct {
	ID int64 `json:"id"`
}
type CreateBodyTypeReq struct {
	NameTM    string   `json:"name_tm"`
	NameEN    string   `json:"name_en"`
	NameRU    string   `json:"name_ru"`
	ImagePath []string `json:"image_path"`
	Category  string   `json:"category"`
	UploadId  string   `json:"upload_id"`
}

type UpdateBodyTypeReq struct {
	ID        int64    `json:"id"`
	NameTM    string   `json:"name_tm"`
	NameEN    string   `json:"name_en"`
	NameRU    string   `json:"name_ru"`
	ImagePath []string `json:"image_path"`
	Category  string   `json:"category"`
	UploadId  string   `json:"upload_id"`
}

type BodyType struct {
	ID        int64    `json:"id"`
	NameTM    string   `json:"name_tm"`
	NameEN    string   `json:"name_en"`
	NameRU    string   `json:"name_ru"`
	ImagePath []string `json:"image_path"`
	Category  string   `json:"category"`
	UploadId  string   `json:"upload_id"`
}

type BodyTypeResult struct {
	BodyTypes []BodyType `json:"body_types"`
	Count     int64      `json:"count"`
}
type CreateBrandReq struct {
	Name       string   `json:"name"`
	LogoPath   []string `json:"logo_path"`
	UploadId   string   `json:"upload_id"`
	Categories []string `json:"categories" binding:"required,dive,oneof=auto moto truck"`
}

type UpdateBrandReq struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	LogoPath   []string `json:"logo_path"`
	UploadId   string   `json:"upload_id"`
	Categories []string `json:"categories" binding:"required,dive,oneof=auto moto truck"`
}
type Brand struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	LogoPath   []string `json:"logo_path"`
	UploadId   string   `json:"upload_id"`
	Categories []string `json:"categories"`
}

type BrandResult struct {
	Brands []Brand `json:"brands"`
	Count  int64   `json:"count"`
}

type CreateModelReq struct {
	Name     string `json:"name"`
	BrandID  int64  `json:"brand_id"`
	Category string `json:"category"`
}

type UpdateModelReq struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	BrandID  int64  `json:"brand_id"`
	Category string `json:"category"`
}

type Model struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	LogoPath  []string `json:"logo_path"`
	UploadId  string   `json:"upload_id"`
	BrandID   int64    `json:"brand_id"`
	BrandName string   `json:"brand_name"`
	Category  string   `json:"category"`
}

type ModelResult struct {
	Models []Model `json:"models"`
	Count  int64   `json:"count"`
}

type CreateDescription struct {
	NameTM string `json:"name_tm"`
	NameEN string `json:"name_en"`
	NameRU string `json:"name_ru"`
}

type Description struct {
	ID     int64  `json:"id"`
	NameTM string `json:"name_tm"`
	NameEN string `json:"name_en"`
	NameRU string `json:"name_ru"`
}

type DescriptionResult struct {
	Descriptions []Description `json:"descriptions"`
	Count        int64         `json:"count"`
}

type UpdateDescription struct {
	ID     int64  `json:"id"`
	NameTM string `json:"name_tm"`
	NameEN string `json:"name_en"`
	NameRU string `json:"name_ru"`
}
