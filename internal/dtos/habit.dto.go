package dtos

import (
	"errors"
	"habitask-backend-go/internal/models"
	"habitask-backend-go/pkg/helpers"
	"strconv"
	"strings"
	"time"
)

type FrequencyDTO struct {
	Kind  string `json:"kind"`            // daily | weekdays | timesPerWeek
	Days  []int  `json:"days,omitempty"`  // weekdays
	Count int    `json:"count,omitempty"` // timesPerWeek
}

type HabitDTO struct {
	ID           string          `json:"id"`
	Type         string          `json:"type"`
	Title        string          `json:"title"`
	Notes        string          `json:"notes"`
	Icon         string          `json:"icon"`
	Color        string          `json:"color"`
	ArchivedAt   *time.Time      `json:"archived_at,omitempty"`
	Frequency    FrequencyDTO    `json:"frequency"`
	StartDate    string          `json:"start_date"`
	TargetPerDay int             `json:"target_per_day"`
	ReminderTime *string         `json:"reminder_time,omitempty"`
	Polarity     models.Polarity `json:"polarity"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	UserID       string          `json:"user_id"`
	User         *UserDTO        `json:"user,omitempty"`
}

func ParseHabit(habit *models.Habit) *HabitDTO {
	if habit == nil {
		return nil
	}
	dto := &HabitDTO{
		ID:           habit.ID,
		Type:         "habit",
		Title:        habit.Title,
		Notes:        habit.Notes,
		Icon:         habit.Icon,
		Color:        habit.Color,
		ArchivedAt:   habit.ArchivedAt,
		Frequency:    columnsToFrequency(habit.FrequencyKind, habit.FrequencyDays, habit.FrequencyCount),
		TargetPerDay: habit.TargetPerDay,
		ReminderTime: habit.ReminderTime,
		Polarity:     habit.Polarity,
		CreatedAt:    habit.CreatedAt,
		UpdatedAt:    habit.UpdatedAt,
		UserID:       habit.UserID,
		User:         ParseUser(habit.User),
	}
	if !habit.StartDate.IsZero() {
		dto.StartDate = habit.StartDate.Format("2006-01-02")
	}
	return dto
}

func ParseHabits(habits []models.Habit) []HabitDTO {
	result := make([]HabitDTO, 0, len(habits))
	for i := range habits {
		result = append(result, *ParseHabit(&habits[i]))
	}
	return result
}

type HabitCreateDTO struct {
	Title        string          `json:"title" binding:"required"`
	Notes        string          `json:"notes"`
	Icon         string          `json:"icon"`
	Color        string          `json:"color"`
	Frequency    FrequencyDTO    `json:"frequency"`
	StartDate    string          `json:"start_date"`
	TargetPerDay int             `json:"target_per_day"`
	ReminderTime *string         `json:"reminder_time"`
	Polarity     models.Polarity `json:"polarity"`
}

func (h *HabitCreateDTO) ToHabit() (*models.Habit, error) {
	kind, days, count := frequencyToColumns(h.Frequency)
	habit := &models.Habit{
		Title:          h.Title,
		Notes:          h.Notes,
		Icon:           h.Icon,
		Color:          h.Color,
		FrequencyKind:  kind,
		FrequencyDays:  days,
		FrequencyCount: count,
		TargetPerDay:   h.TargetPerDay,
		ReminderTime:   h.ReminderTime,
		Polarity:       h.Polarity,
	}
	if habit.TargetPerDay < 1 {
		habit.TargetPerDay = 1
	}
	if habit.Polarity == "" {
		habit.Polarity = models.PolarityPositive
	}
	if h.StartDate != "" {
		sd, err := helpers.ParseDate(h.StartDate)
		if err != nil {
			return nil, errors.New("invalid date format on field start_date")
		}
		habit.StartDate = sd
	} else {
		habit.StartDate = time.Now()
	}
	return habit, nil
}

type HabitUpdateDTO struct {
	Title        *string          `json:"title"`
	Notes        *string          `json:"notes"`
	Icon         *string          `json:"icon"`
	Color        *string          `json:"color"`
	Frequency    *FrequencyDTO    `json:"frequency"`
	StartDate    *string          `json:"start_date"`
	TargetPerDay *int             `json:"target_per_day"`
	ReminderTime *string          `json:"reminder_time"`
	Polarity     *models.Polarity `json:"polarity"`
	Archived     *bool            `json:"archived"`
}

func (h *HabitUpdateDTO) ToHabit() (*models.Habit, error) {
	habit := &models.Habit{}
	if h.Title != nil {
		habit.Title = *h.Title
	}
	if h.Notes != nil {
		habit.Notes = *h.Notes
	}
	if h.Icon != nil {
		habit.Icon = *h.Icon
	}
	if h.Color != nil {
		habit.Color = *h.Color
	}
	if h.Frequency != nil {
		kind, days, count := frequencyToColumns(*h.Frequency)
		habit.FrequencyKind = kind
		habit.FrequencyDays = days
		habit.FrequencyCount = count
	}
	if h.StartDate != nil {
		sd, err := helpers.ParseDate(*h.StartDate)
		if err != nil {
			return nil, errors.New("invalid date format on field start_date")
		}
		habit.StartDate = sd
	}
	if h.TargetPerDay != nil {
		habit.TargetPerDay = *h.TargetPerDay
	}
	if h.ReminderTime != nil {
		habit.ReminderTime = h.ReminderTime
	}
	if h.Polarity != nil {
		habit.Polarity = *h.Polarity
	}
	if h.Archived != nil {
		if *h.Archived {
			now := time.Now()
			habit.ArchivedAt = &now
		} else {
			zero := time.Time{}
			habit.ArchivedAt = &zero
		}
	}
	return habit, nil
}

func frequencyToColumns(f FrequencyDTO) (models.FrequencyKind, *string, *int) {
	kind := models.FrequencyKind(f.Kind)
	if kind == "" {
		kind = models.FrequencyDaily
	}
	var days *string
	var count *int
	if kind == models.FrequencyWeekdays && len(f.Days) > 0 {
		parts := make([]string, 0, len(f.Days))
		for _, d := range f.Days {
			parts = append(parts, strconv.Itoa(d))
		}
		joined := strings.Join(parts, ",")
		days = &joined
	}
	if kind == models.FrequencyTimesPerWeek && f.Count > 0 {
		c := f.Count
		count = &c
	}
	return kind, days, count
}

func columnsToFrequency(kind models.FrequencyKind, days *string, count *int) FrequencyDTO {
	if kind == "" {
		kind = models.FrequencyDaily
	}
	f := FrequencyDTO{Kind: string(kind)}
	if days != nil && *days != "" {
		for _, p := range strings.Split(*days, ",") {
			if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
				f.Days = append(f.Days, n)
			}
		}
	}
	if count != nil {
		f.Count = *count
	}
	return f
}
