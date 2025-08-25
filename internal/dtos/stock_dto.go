package dtos

type CreateStockReq struct {
	UserID      int64    `json:"user_id"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
	StoreName   string   `json:"store_name"`
	Images      []string `json:"images"`
	Logo        string   `json:"logo"`
	RegionID    int64    `json:"region_id"`
	CityID      int64    `json:"city_id"`
	Address     string   `json:"address"`
	Status      string   `json:"status"`
	Description string   `json:"description"`
}

type UpdateStockReq struct {
	ID          int64    `json:"id"`
	UserID      int64    `json:"user_id"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
	StoreName   string   `json:"store_name"`
	Images      []string `json:"images"`
	Logo        string   `json:"logo"`
	RegionID    int64    `json:"region_id"`
	CityID      int64    `json:"city_id"`
	Address     string   `json:"address"`
	Status      string   `json:"status"`
	Description string   `json:"description"`
}

type Stock struct {
	ID           int64    `json:"id"`
	UserID       int64    `json:"user_id"`
	UserName     string   `json:"user_name"`
	PhoneNumber  string   `json:"phone_number"`
	Email        string   `json:"email"`
	StoreName    string   `json:"store_name"`
	Images       []string `json:"images"`
	Logo         string   `json:"logo"`
	RegionID     int64    `json:"region_id"`
	CityID       int64    `json:"city_id"`
	Address      string   `json:"address"`
	CityNameTM   string   `json:"city_name_tm"`
	CityNameEN   string   `json:"city_name_en"`
	CityNameRU   string   `json:"city_name_ru"`
	RegionNameTM string   `json:"region_name_tm"`
	RegionNameEN string   `json:"region_name_en"`
	RegionNameRU string   `json:"region_name_ru"`
	Status       string   `json:"status"`
	Description  string   `json:"description"`
}

type StocksResult struct {
	Stocks []Stock `json:"stocks"`
	Count  int64   `json:"count"`
}

type UpdateStockStatus struct {
	ID      int64  `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ReqSendPushDTO struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
