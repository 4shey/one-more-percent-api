package models

// Schedule represents a recurring daily schedule entry.
type Schedule struct {
	ID        int
	UserID    int
	DayOfWeek string // "Monday", "Tuesday", etc.
	StartTime string // "HH:MM"
	EndTime   string // "HH:MM"
	Activity  string
	IsActive  bool
}
