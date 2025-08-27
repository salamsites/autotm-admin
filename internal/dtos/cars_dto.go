package dtos

type Car struct {
	Id             int64       `json:"id"`
	UserId         int64       `json:"user_id"`
	UserName       string      `json:"user_name"`
	StockId        int64       `json:"stock_id"`
	StoreName      string      `json:"store_name"`
	BrandId        int64       `json:"brand_id"`
	BrandName      string      `json:"brand_name"`
	ModelId        int64       `json:"model_id"`
	ModelName      string      `json:"model_name"`
	Year           int64       `json:"year"`
	Mileage        int64       `json:"mileage"`
	Color          string      `json:"color"`
	EngineCapacity float64     `json:"engine_capacity"`
	EngineType     string      `json:"engine_type"`
	BodyId         int64       `json:"body_id"`
	BodyNameTM     string      `json:"body_name_tm"`
	BodyNameEN     string      `json:"body_name_en"`
	BodyNameRU     string      `json:"body_name_ru"`
	Transmission   string      `json:"transmission"`
	DriveType      string      `json:"drive_type"`
	Vin            string      `json:"vin"`
	Description    string      `json:"description"`
	CityId         int64       `json:"city_id"`
	CityNameTM     string      `json:"city_name_tm"`
	CityNameEN     string      `json:"city_name_en"`
	CityNameRU     string      `json:"city_name_ru"`
	Name           string      `json:"name"`
	Mail           string      `json:"mail"`
	PhoneNumber    string      `json:"phone_number"`
	Price          int64       `json:"price"`
	IsComment      bool        `json:"is_comment"`
	IsExchange     bool        `json:"is_exchange"`
	IsCredit       bool        `json:"is_credit"`
	Images         interface{} `json:"images"`
	Status         string      `json:"status"`
}

type CarsResp struct {
	Cars  []Car `json:"cars"`
	Count int64 `json:"count"`
}

type UpdateCarStatus struct {
	ID      int64  `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Truck struct {
	Id              int64       `json:"id"`
	UserId          int64       `json:"user_id"`
	UserName        string      `json:"user_name"`
	StockId         int64       `json:"stock_id"`
	StoreName       string      `json:"store_name"`
	BodyId          int64       `json:"body_id"`
	BodyNameTM      string      `json:"body_name_tm"`
	BodyNameEN      string      `json:"body_name_en"`
	BodyNameRU      string      `json:"body_name_ru"`
	BrandId         int64       `json:"brand_id"`
	BrandName       string      `json:"brand_name"`
	ModelId         int64       `json:"model_id"`
	ModelName       string      `json:"model_name"`
	LoadCapacity    float64     `json:"load_capacity"`
	Price           int64       `json:"price"`
	BodyType        string      `json:"body_type"`
	DriveType       string      `json:"drive_type"`
	Transmission    string      `json:"transmission"`
	EngineType      string      `json:"engine_type"`
	Year            int64       `json:"year"`
	Seats           int64       `json:"seats"`
	CabType         string      `json:"cab_type"`
	WheelFormula    string      `json:"wheel_formula"`
	Chassis         string      `json:"chassis"`
	CubSuspension   string      `json:"cub_suspension"`
	BusType         string      `json:"bus_type"`
	SuspensionType  string      `json:"suspension_type"`
	Brakes          string      `json:"brakes"`
	Axles           int64       `json:"axles"`
	EngineHours     int64       `json:"engine_hours"`
	VehicleType     string      `json:"vehicle_type"`
	EngineCapacity  float64     `json:"engine_capacity"`
	ForkLiftType    string      `json:"forklift_type"`
	LiftingCapacity float64     `json:"lifting_capacity"`
	Mileage         int64       `json:"mileage"`
	ExcavatorType   string      `json:"excavator_type"`
	BulldozerType   string      `json:"bulldozer_type"`
	Color           string      `json:"color"`
	Vin             string      `json:"vin"`
	Description     string      `json:"description"`
	CityId          int64       `json:"city_id"`
	CityNameTM      string      `json:"city_name_tm"`
	CityNameEN      string      `json:"city_name_en"`
	CityNameRU      string      `json:"city_name_ru"`
	Name            string      `json:"name"`
	Mail            string      `json:"mail"`
	PhoneNumber     string      `json:"phone_number"`
	IsComment       bool        `json:"is_comment"`
	IsExchange      bool        `json:"is_exchange"`
	IsCredit        bool        `json:"is_credit"`
	Images          interface{} `json:"images"`
	Status          string      `json:"status"`
}

type TrucksResp struct {
	Trucks []Truck `json:"trucks"`
	Count  int64   `json:"count"`
}

type UpdateTruckStatus struct {
	ID      int64  `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
