package models

type AutoStore struct {
	ID          int64
	UserID      int64
	PhoneNumber string
	Email       string
	StoreName   string
	Images      []string
	LogoPath    string
	Address     string
	RegionID    int64
	CityID      int64
	CityName    string
	RegionName  string
	UserName    string
}
