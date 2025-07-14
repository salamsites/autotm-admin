package dtos

type CreateRegionReq struct {
	NameTM string `json:"name_tm"`
	NameEN string `json:"name_en"`
	NameRu string `json:"name_ru"`
}

type UpdateRegionReq struct {
	ID     int64  `json:"id"`
	NameTM string `json:"name_tm"`
	NameEN string `json:"name_en"`
	NameRu string `json:"name_ru"`
}
type Region struct {
	ID     int64  `json:"id"`
	NameTM string `json:"name_tm"`
	NameEN string `json:"name_en"`
	NameRu string `json:"name_ru"`
}

type RegionResult struct {
	Regions []Region `json:"regions"`
	Count   int64    `json:"count"`
}

type CreateCityReq struct {
	NameTM   string `json:"name_tm"`
	NameEN   string `json:"name_en"`
	NameRu   string `json:"name_ru"`
	RegionID int64  `json:"region_id"`
}

type UpdateCityReq struct {
	ID       int64  `json:"id"`
	NameTM   string `json:"name_tm"`
	NameEN   string `json:"name_en"`
	NameRu   string `json:"name_ru"`
	RegionID int64  `json:"region_id"`
}

type City struct {
	ID           int64  `json:"id"`
	NameTM       string `json:"name_tm"`
	NameEN       string `json:"name_en"`
	NameRu       string `json:"name_ru"`
	RegionID     int64  `json:"region_id"`
	RegionNameTM string `json:"region_name_tm"`
	RegionNameEN string `json:"region_name_en"`
	RegionNameRU string `json:"region_name_ru"`
}

type CityResult struct {
	Cities []City `json:"cities"`
	Count  int64  `json:"count"`
}
