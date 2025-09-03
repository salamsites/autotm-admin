package models

import "database/sql"

type GetUser struct {
	Id          int64
	FullName    sql.NullString
	Email       sql.NullString
	PhoneNumber sql.NullString
}
