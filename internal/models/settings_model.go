package models

import "encoding/json"

type Role struct {
	ID   int64
	Name string
	Role json.RawMessage
}

type User struct {
	ID       int64
	Username string
	Login    string
	Password string
	RoleID   int64
	RoleName string
}
