package models

type Reminder struct {
	BaseModel
	Title               string  `gorm:"unique;varchar(100)"`
	Description         string  `gorm:"varchar(500)"`
	Purpose             string  `gorm:"varchar(255)"`
	Category            string  `gorm:"enum('alarm', 'reminder');Varchar(50)"`
	Trigger             string  `gorm:"enum('datetime', 'time', 'interval','source');Varchar(50)"`
	TriggerValue        string  `gorm:"varchar(50)"`
	TriggerIntervalUnit *string `gorm:"enum('seconds', 'minutes', 'hours', 'days', 'weeks', 'months', 'years');Varchar(50)"`
	Source              *string `gorm:"enum('habit', 'habit-action', 'task', 'task_step');Varchar(50)"`
	SourceID            *string `gorm:"index"`
	SourceField         *string `gorm:"varchar(50)"`
	SourceFieldType     *string `gorm:"enum('string', 'number', 'date', 'time');Varchar(50)" json:"-"`
	UserID              string  `gorm:"index;type:varchar(36)"`
	User                *User   `gorm:"foreignKey:UserID"`
	IsRead              *bool   `gorm:"default:false"`
}

type ReminderQuery struct {
	Reminder
	Search  string  `json:"search"`
	DateLte *string `json:"date_lte"`
	DateGte *string `json:"date_gte"`
}
