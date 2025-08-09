package dtos

type GetUser struct {
	Id          int64   `json:"id"`
	FullName    string  `json:"full_name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

type GetUsersResult struct {
	GetUsers []GetUser `json:"get_users"`
	Count    int64     `json:"count"`
}
