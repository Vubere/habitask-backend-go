package models

type DeviceToken struct {
	BaseModel
	UserID   string `gorm:"type:varchar(36);index;not null" json:"user_id"`
	Token    string `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	Platform string `gorm:"type:varchar(20)" json:"platform"`
}
