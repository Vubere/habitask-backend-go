package models

import "time"

type FrequencyKind string

const (
	FrequencyDaily        FrequencyKind = "daily"
	FrequencyWeekdays     FrequencyKind = "weekdays"
	FrequencyTimesPerWeek FrequencyKind = "timesPerWeek"
)

type Polarity string

const (
	PolarityPositive Polarity = "positive"
	PolarityNegative Polarity = "negative"
)

type Habit struct {
	BaseModel
	Title          string        `gorm:"varchar(100);not null" json:"title"`
	Notes          string        `gorm:"varchar(500)" json:"notes"`
	Icon           string        `gorm:"varchar(50)" json:"icon"`
	Color          string        `gorm:"varchar(20)" json:"color"`
	ArchivedAt     *time.Time    `gorm:"datetime" json:"archived_at,omitempty"`
	UserID         string        `gorm:"index;varchar(36)" json:"user_id"`
	User           *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	FrequencyKind  FrequencyKind `gorm:"varchar(20);not null;default:daily" json:"frequency_kind"`
	FrequencyDays  *string       `gorm:"varchar(20)" json:"frequency_days,omitempty"` // "1,3,5"
	FrequencyCount *int          `json:"frequency_count,omitempty"`                   // Per week count
	StartDate      time.Time     `gorm:"type:date;not null" json:"start_date"`
	TargetPerDay   int           `gorm:"default:1" json:"target_per_day"`
	ReminderTime   *string       `gorm:"varchar(5)" json:"reminder_time,omitempty"` // "HH:mm"
	Polarity       Polarity      `gorm:"varchar(10);default:positive" json:"polarity"`
}

type HabitSummary struct {
	Label         string  `json:"label"`
	Count         int     `json:"count"`
	PositiveCount int     `json:"positive_count"`
	NegativeCount int     `json:"negative_count"`
	TargetTotal   int     `json:"target_total"`
	LastHabitDate *string `json:"last_habit_date"`
}

type HabitQuery struct {
	Habit
	Search   string     `json:"search"`
	DateLte  *time.Time `json:"date_lte"`
	DateGte  *time.Time `json:"date_gte"`
	Archived *bool      `json:"archived"`
	//SummaryFilters
	GroupBy   string `json:"group_by"`
	DateGroup string `json:"date_group"` //DAY, WEEK, MONTH, YEAR
}

type HabitStats struct {
	CurrentStreak  int     `json:"current_streak"`
	BestStreak     int     `json:"best_streak"`
	CompletionRate float64 `json:"completion_rate"`
}
