package models

type BodyType struct {
	ID        int64
	Name      string
	ImagePath string
}
type Brand struct {
	ID         int64
	Name       string
	LogoPath   string
	Categories []string
}

type ID struct {
	ID       int64  `json:"id"`
	Category string `json:"category"`
}

type Model struct {
	ID           int64
	Name         string
	LogoPath     string
	BrandID      int64
	BrandName    string
	BodyTypeID   int64
	BodyTypeName string
}
