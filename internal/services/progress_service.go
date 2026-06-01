package services

import (
	"database/sql"
	"fmt"
	"one_more_percent/internal/db"
	"one_more_percent/internal/models"
	"time"
)

// GetSchedulesForDay returns all active schedules for a given day name (e.g. "Monday").
func GetSchedulesForDay(day string) ([]models.Schedule, error) {
	rows, err := db.DB.Query(`
		SELECT id, user_id, day_of_week,
		       TO_CHAR(start_time, 'HH24:MI'),
		       TO_CHAR(end_time, 'HH24:MI'),
		       activity, is_active
		FROM schedules
		WHERE day_of_week = $1 AND is_active = TRUE
		ORDER BY start_time
	`, day)
	if err != nil {
		return nil, fmt.Errorf("query schedules: %w", err)
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.DayOfWeek,
			&s.StartTime, &s.EndTime, &s.Activity, &s.IsActive,
		); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

// EnsureProgressRow creates a pending progress row if it doesn't exist yet.
func EnsureProgressRow(scheduleID int, date time.Time) error {
	_, err := db.DB.Exec(`
		INSERT INTO schedule_progresses (schedule_id, progress_date, status)
		VALUES ($1, $2, 'pending')
		ON CONFLICT (schedule_id, progress_date) DO NOTHING
	`, scheduleID, date.Format("2006-01-02"))
	return err
}

// MarkCompleted sets a pending progress row to completed.
func MarkCompleted(scheduleID int, date time.Time) error {
	_, err := db.DB.Exec(`
		UPDATE schedule_progresses
		SET status = 'completed', completed_at = NOW(), updated_at = NOW()
		WHERE schedule_id = $1 AND progress_date = $2 AND status = 'pending'
	`, scheduleID, date.Format("2006-01-02"))
	return err
}

// MarkAllPendingAsMissed flips all pending rows on a date to missed.
func MarkAllPendingAsMissed(date time.Time) error {
	_, err := db.DB.Exec(`
		UPDATE schedule_progresses
		SET status = 'missed', updated_at = NOW()
		WHERE progress_date = $1 AND status = 'pending'
	`, date.Format("2006-01-02"))
	return err
}

// GetDayProgress returns completed and missed progress entries for a given date.
func GetDayProgress(date time.Time) (completed []models.ScheduleProgress, missed []models.ScheduleProgress, err error) {
	rows, err := db.DB.Query(`
		SELECT p.id, p.schedule_id, p.progress_date, p.status,
		       s.activity,
		       TO_CHAR(s.start_time, 'HH24:MI'),
		       TO_CHAR(s.end_time, 'HH24:MI')
		FROM schedule_progresses p
		JOIN schedules s ON s.id = p.schedule_id
		WHERE p.progress_date = $1
		ORDER BY s.start_time
	`, date.Format("2006-01-02"))
	if err != nil {
		return nil, nil, fmt.Errorf("query progress: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.ScheduleProgress
		if err := rows.Scan(
			&p.ID, &p.ScheduleID, &p.ProgressDate, &p.Status,
			&p.Activity, &p.StartTime, &p.EndTime,
		); err != nil {
			return nil, nil, err
		}
		if p.Status == "completed" {
			completed = append(completed, p)
		} else {
			missed = append(missed, p)
		}
	}
	return completed, missed, nil
}

// GetActiveSchedule returns the schedule currently in progress (start_time <= now < end_time).
// Returns nil if no schedule is active right now.
func GetActiveSchedule(dayName, currentHHMM string) (*models.Schedule, error) {
	var s models.Schedule
	err := db.DB.QueryRow(`
		SELECT id, user_id, day_of_week,
		       TO_CHAR(start_time, 'HH24:MI'),
		       TO_CHAR(end_time, 'HH24:MI'),
		       activity, is_active
		FROM schedules
		WHERE day_of_week = $1
		  AND start_time <= $2::time
		  AND end_time   >  $2::time
		  AND is_active = TRUE
		ORDER BY start_time
		LIMIT 1
	`, dayName, currentHHMM).Scan(
		&s.ID, &s.UserID, &s.DayOfWeek,
		&s.StartTime, &s.EndTime, &s.Activity, &s.IsActive,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetActiveSchedule: %w", err)
	}
	return &s, nil
}

// ProgressRowExists checks whether a progress row already exists for a schedule on a given date.
// Used by the catch-up check to avoid double-sending reminders.
func ProgressRowExists(scheduleID int, date time.Time) (bool, error) {
	var count int
	err := db.DB.QueryRow(`
		SELECT COUNT(*) FROM schedule_progresses
		WHERE schedule_id = $1 AND progress_date = $2
	`, scheduleID, date.Format("2006-01-02")).Scan(&count)
	return count > 0, err
}
