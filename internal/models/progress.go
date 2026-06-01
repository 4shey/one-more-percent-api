package models

import "time"

// ScheduleProgress tracks the completion status of a schedule on a specific date.
type ScheduleProgress struct {
	ID           int
	ScheduleID   int
	ProgressDate time.Time
	Status       string // pending | completed | missed
	CompletedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Joined from schedules
	Activity  string
	StartTime string
	EndTime   string
}
