package models

import "time"

type User struct {
	BaseModel
	FirstName   string     `gorm:"Varchar(50)"`
	LastName    string     `gorm:"Varchar(50)"`
	Username    string     `gorm:"unique;Varchar(100);"`
	Email       string     `gorm:"unique;Varchar(100)"`
	Password    string     `gorm:"Varchar(100)" json:"-"`
	Token       string     `gorm:"unique;Varchar(100)"`
	TokenExpiry *time.Time `gorm:"datetime"`
	IsAdmin     bool       `gorm:"default:false"`
	IsActive    bool       `gorm:"default:false"`
	Profession  string     `gorm:"Varchar(50)"`
	Bio         string     `gorm:"Varchar(500)"`
}

type UserQuery struct {
	User
	Search string `json:"search"`
	Active *bool  `json:"active"`
	Admin  *bool  `json:"admin"`
}
