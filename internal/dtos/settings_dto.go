package dtos

import "encoding/json"

type Role struct {
	ID   int64           `json:"id"`
	Name string          `json:"name"`
	Role json.RawMessage `json:"role"`
}

type RoleResult struct {
	Roles []Role `json:"roles"`
	Count int64  `json:"count"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username" validate:"required"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	RoleID   int64  `json:"role_id" validate:"required"`
	RoleName string `json:"role_name"`
}

type UserResult struct {
	Users []User `json:"users"`
	Count int64  `json:"count"`
}
