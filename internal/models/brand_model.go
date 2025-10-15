package models

type BodyType struct {
	ID        int64
	NameTM    string
	NameEN    string
	NameRU    string
	ImagePath interface{}
	Category  string
	UploadId  string
}
type Brand struct {
	ID         int64
	Name       string
	LogoPath   interface{}
	UploadId   string
	Categories []string
}

type Model struct {
	ID        int64
	Name      string
	LogoPath  interface{}
	UploadId  string
	BrandID   int64
	BrandName string
	Category  string
}

type Description struct {
	ID         int64
	NameTM     string
	NameEN     string
	NameRU     string
	Categories []string
}
