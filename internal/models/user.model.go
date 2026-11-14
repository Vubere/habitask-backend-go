package models

import "time"

type User struct {
	BaseModel
	FirstName   string     `gorm:"varchar(50)"`
	LastName    string     `gorm:"varchar(50)"`
	Username    string     `gorm:"unique;varchar(100);"`
	Email       string     `gorm:"unique;varchar(100)"`
	Password    string     `gorm:"varchar(100)" json:"-"`
	Token       string     `gorm:"unique;varchar(100)"`
	TokenExpiry *time.Time `gorm:"datetime"`
	IsAdmin     bool       `gorm:"default:false"`
	IsActive    bool       `gorm:"default:false"`
	Profession  string     `gorm:"varchar(50)"`
	Bio         string     `gorm:"varchar(500)"`
}

type UserQuery struct {
	User
	Search string `json:"search"`
	Active *bool  `json:"active"`
	Admin  *bool  `json:"admin"`
}
