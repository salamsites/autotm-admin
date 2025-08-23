package models

type Stock struct {
	ID           int64
	UserID       int64
	UserName     string
	PhoneNumber  string
	Email        string
	StoreName    string
	Images       []string
	Logo         string
	Address      string
	RegionID     int64
	CityID       int64
	CityNameTM   string
	CityNameEN   string
	CityNameRU   string
	RegionNameTM string
	RegionNameEN string
	RegionNameRU string
	Status       string
	Description  string
}
