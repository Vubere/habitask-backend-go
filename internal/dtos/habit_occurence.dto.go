package dtos

type HabitOccurenceSummaryDTO struct {
	HabitID    string    `json:"habit_id"`
	Habbit     *HabitDTO `json:"habit,omitempty"`
	Occurences int       `json:"occurences"`
	LastDone   *string   `json:"last_done"`
}
