package models

import "time"

type AdminUser struct {
	Id int
	UserName string
	Password string
	Desc string
	Status int
	CreateTime time.Time
}

func (AdminUser) TableName() string {
	return "admin_users"

}
