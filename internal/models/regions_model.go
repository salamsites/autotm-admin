package models

type Region struct {
	ID     int64
	NameTM string
	NameEN string
	NameRU string
}

type City struct {
	ID           int64
	NameTM       string
	NameEN       string
	NameRU       string
	RegionID     int64
	RegionNameTM string
	RegionNameEN string
	RegionNameRU string
}
