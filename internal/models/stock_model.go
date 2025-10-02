package models

type Stock struct {
	ID           int64
	UserID       int64
	UserName     *string
	PhoneNumber  string
	Email        string
	StoreName    string
	Images       interface{}
	Logo         interface{}
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
	Location     interface{}
}

type ESStock struct {
	ID           int64       `json:"id" es:"type=keyword"`
	UserID       int64       `json:"user_id" es:"type=long"`
	UserName     string      `json:"user_name" es:"type=text"`
	PhoneNumber  string      `json:"phone_number" es:"type=keyword"`
	Email        string      `json:"email" es:"type=keyword"`
	StoreName    string      `json:"store_name" es:"type=text"`
	Images       interface{} `json:"images" es:"type=object"`
	Logo         interface{} `json:"logo" es:"type=object"`
	Address      string      `json:"address" es:"type=text"`
	RegionID     int64       `json:"region_id" es:"type=long"`
	CityID       int64       `json:"city_id" es:"type=long"`
	CityNameTM   string      `json:"city_name_tm" es:"type=text"`
	CityNameEN   string      `json:"city_name_en" es:"type=text"`
	CityNameRU   string      `json:"city_name_ru" es:"type=text"`
	RegionNameTM string      `json:"region_name_tm" es:"type=text"`
	RegionNameEN string      `json:"region_name_en" es:"type=text"`
	RegionNameRU string      `json:"region_name_ru" es:"type=text"`
	Status       string      `json:"status" es:"type=keyword"`
	Description  string      `json:"description" es:"type=text"`
	Location     interface{} `json:"location" es:"type=object"`
}
