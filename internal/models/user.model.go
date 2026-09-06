package models

import "time"

type User struct {
	BaseModel
	FirstName   string     `gorm:"Varchar(50)"`
	LastName    string     `gorm:"Varchar(50)"`
	Email       string     `gorm:"unique;Varchar(100)"`
	Password    string     `gorm:"Varchar(100)"`
	Token       string     `gorm:"unique;Varchar(100)"`
	TokenExpiry *time.Time `gorm:"datetime"`
	IsAdmin     bool       `gorm:"default:false"`
	IsActive    bool       `gorm:"default:false"`
	Profession  string     `gorm:"Varchar(50)"`
	Bio         string     `gorm:"Varchar(500)"`
}
