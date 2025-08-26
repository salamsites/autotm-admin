package models

type Role struct {
	ID   int64
	Name string
	Role interface{}
}

type User struct {
	ID       int64
	Username string
	Login    string
	Password string
	RoleID   int64
	RoleName string
	Status   bool
}
