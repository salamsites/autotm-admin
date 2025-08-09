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
}

type Stock struct {
	ID           int64    `json:"id"`
	UserID       int64    `json:"user_id"`
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
	UserName     string   `json:"user_name"`
}

type StocksResult struct {
	Stocks []Stock `json:"stocks"`
	Count  int64   `json:"count"`
}

type GetUsers struct {
	Id          int64   `json:"id"`
	FullName    string  `json:"full_name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

type GetUserResult struct {
	Users []GetUsers `json:"users"`
	Count int64      `json:"count"`
}

type GetUserByIDsReq struct {
	Ids []int64 `json:"ids"`
}
