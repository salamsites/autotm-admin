package dtos

type CreateAutoStoreReq struct {
	UserID      int64    `json:"user_id"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
	StoreName   string   `json:"store_name" binding:"required"`
	Images      []string `json:"images"`
	LogoPath    string   `json:"logo_path"`
	RegionID    int64    `json:"region_id"`
	CityID      int64    `json:"city_id"`
	Address     string   `json:"address"`
}

type UpdateAutoStoreReq struct {
	ID          int64    `json:"id"`
	UserID      int64    `json:"user_id"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
	StoreName   string   `json:"store_name" binding:"required"`
	Images      []string `json:"images"`
	LogoPath    string   `json:"logo_path"`
	RegionID    int64    `json:"region_id"`
	CityID      int64    `json:"city_id"`
	Address     string   `json:"address"`
}

type AutoStore struct {
	ID          int64    `json:"id"`
	UserID      int64    `json:"user_id"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
	StoreName   string   `json:"store_name"`
	Images      []string `json:"images"`
	LogoPath    string   `json:"logo_path"`
	RegionID    int64    `json:"region_id"`
	CityID      int64    `json:"city_id"`
	Address     string   `json:"address"`
	CityName    string   `json:"city_name"`
	RegionName  string   `json:"region_name"`
	UserName    *string  `json:"user_name"`
}

type AutoStoresResult struct {
	AutoStores []AutoStore `json:"auto_stores"`
	Count      int64       `json:"count"`
}

type GetUsers struct {
	Id          int64   `json:"id"`
	FullName    *string `json:"full_name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Avatar      *string `json:"avatar"`
}

type GetUserResult struct {
	Users []GetUsers `json:"users"`
	Count int64      `json:"count"`
}

type GetUserByIDsReq struct {
	Ids []int64 `json:"ids"`
}
