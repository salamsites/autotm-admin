package dtos

type CreateRoleReq struct {
	Name string      `json:"name"`
	Role interface{} `json:"role"`
}

type UpdateRoleReq struct {
	ID   int64       `json:"id"`
	Name string      `json:"name"`
	Role interface{} `json:"role"`
}

type Role struct {
	ID   int64       `json:"id"`
	Name string      `json:"name"`
	Role interface{} `json:"role"`
}

type RoleResult struct {
	Roles []Role `json:"roles"`
	Count int64  `json:"count"`
}

type CreateUserReq struct {
	Username string `json:"username" validate:"required"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	RoleID   int64  `json:"role_id" validate:"required"`
	Status   bool   `json:"status" validate:"required"`
}

type UpdateUserReq struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	RoleID   int64  `json:"role_id"`
	Status   bool   `json:"status"`
}
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username" validate:"required"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	RoleID   int64  `json:"role_id" validate:"required"`
	RoleName string `json:"role_name"`
	Status   bool   `json:"status"`
}

type UserResult struct {
	Users []User `json:"users"`
	Count int64  `json:"count"`
}

type LoginReq struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResp struct {
	Token string `json:"token"`
}
