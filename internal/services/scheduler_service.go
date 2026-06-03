package services

import (
	"fmt"
	"one_more_percent/internal/models"
	"os"
	"strconv"
	"sync"
	"time"
)

// PendingReminder holds state about the last reminder sent to a user,
// so we know which schedule to mark complete when they reply.
type PendingReminder struct {
	ScheduleID   int
	Activity     string
	ProgressDate time.Time
}

var (
	pendingMu      sync.Mutex
	pendingReminders = map[int64]*PendingReminder{}
)

func SetPendingReminder(chatID int64, r *PendingReminder) {
	pendingMu.Lock()
	defer pendingMu.Unlock()
	pendingReminders[chatID] = r
}

func GetPendingReminder(chatID int64) *PendingReminder {
	pendingMu.Lock()
	defer pendingMu.Unlock()
	return pendingReminders[chatID]
}

func ClearPendingReminder(chatID int64) {
	pendingMu.Lock()
	defer pendingMu.Unlock()
	delete(pendingReminders, chatID)
}

var jakartaLoc *time.Location

// StartScheduler launches the background ticker that drives reminders and midnight recap.
func StartScheduler() {
	var err error
	jakartaLoc, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println("[Scheduler] Failed to load Asia/Jakarta, using UTC:", err)
		jakartaLoc = time.UTC
	}

	// Send Startup Message
	if chatID, err := parseChatID(); err == nil {
		msg := "Woi gw udah nyala nih (baru di-deploy). Coba cek jadwal lu sekarang ada apa ngga. Jangan males-malesan!"
		_ = SendTelegramMessage(chatID, msg)
		fmt.Println("[Scheduler] Startup greeting sent")
	}

	go func() {
		// Align to the next minute boundary so checks fire on :00 seconds.
		now := time.Now().In(jakartaLoc)
		sleepDuration := time.Duration(60-now.Second())*time.Second -
			time.Duration(now.Nanosecond())*time.Nanosecond
		fmt.Printf("[Scheduler] Starting. Next check in %s\n", sleepDuration.Round(time.Second))

		// Catch-up: send reminder immediately if we're already inside a schedule window.
		runCatchUpCheck()

		time.Sleep(sleepDuration)

		// Run once right after alignment, then every minute.
		runCheck()
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			runCheck()
		}
	}()
}

func runCheck() {
	now := time.Now().In(jakartaLoc)
	fmt.Printf("[Scheduler] tick %s\n", now.Format("2006-01-02 15:04"))

	// Midnight recap: at exactly 00:00 recap yesterday.
	if now.Hour() == 0 && now.Minute() == 0 {
		yesterday := now.AddDate(0, 0, -1)
		runMidnightRecap(yesterday)
	}

	// Schedule reminders: find any schedule whose start_time matches now.
	dayName := now.Weekday().String() // e.g. "Monday"
	schedules, err := GetSchedulesForDay(dayName)
	if err != nil {
		fmt.Println("[Scheduler] Error fetching schedules:", err)
		return
	}

	currentHHMM := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())
	for _, s := range schedules {
		if s.StartTime == currentHHMM {
			sendReminder(s, now)
		}
	}
}

func sendReminder(s models.Schedule, now time.Time) {
	chatID, err := parseChatID()
	if err != nil {
		fmt.Println("[Scheduler] Invalid TELEGRAM_CHAT_ID:", err)
		return
	}

	today := truncateToDay(now)

	if err := EnsureProgressRow(s.ID, today); err != nil {
		fmt.Println("[Scheduler] EnsureProgressRow error:", err)
	}

	SetPendingReminder(chatID, &PendingReminder{
		ScheduleID:   s.ID,
		Activity:     s.Activity,
		ProgressDate: today,
	})

	msg := GenerateReminder(s.Activity)
	if err := SendTelegramMessage(chatID, msg); err != nil {
		fmt.Println("[Scheduler] Send reminder error:", err)
	} else {
		fmt.Printf("[Scheduler] Reminder sent for '%s'\n", s.Activity)
	}
}

// runMidnightRecap collects yesterday's progress, marks pending as missed, and sends AI recap.
func runMidnightRecap(date time.Time) {
	fmt.Printf("[Recap] Running midnight recap for %s\n", date.Format("2006-01-02"))

	chatID, err := parseChatID()
	if err != nil {
		fmt.Println("[Recap] Invalid TELEGRAM_CHAT_ID:", err)
		return
	}

	// Mark all remaining pending rows as missed.
	if err := MarkAllPendingAsMissed(date); err != nil {
		fmt.Println("[Recap] MarkAllPendingAsMissed error:", err)
	}

	completed, missed, err := GetDayProgress(date)
	if err != nil {
		fmt.Println("[Recap] GetDayProgress error:", err)
		return
	}

	if len(completed) == 0 && len(missed) == 0 {
		fmt.Println("[Recap] No schedules found for", date.Format("2006-01-02"))
		return
	}

	recap := GenerateRecap(completed, missed)
	if err := SendTelegramMessage(chatID, recap); err != nil {
		fmt.Println("[Recap] Send error:", err)
	} else {
		fmt.Println("[Recap] Recap sent ✓")
	}
}

// --- helpers ---

// runCatchUpCheck fires once at startup. If we're currently inside a schedule window
// and no progress row exists yet (reminder never sent), it sends the reminder now.
func runCatchUpCheck() {
	now := time.Now().In(jakartaLoc)
	dayName := now.Weekday().String()
	currentHHMM := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())

	schedule, err := GetActiveSchedule(dayName, currentHHMM)
	if err != nil {
		fmt.Println("[CatchUp] Error fetching active schedule:", err)
		return
	}
	if schedule == nil {
		fmt.Println("[CatchUp] No active schedule right now, nothing to catch up")
		return
	}

	today := truncateToDay(now)
	exists, err := ProgressRowExists(schedule.ID, today)
	if err != nil {
		fmt.Println("[CatchUp] ProgressRowExists error:", err)
		return
	}
	if exists {
		fmt.Printf("[CatchUp] '%s' already reminded today, skipping\n", schedule.Activity)
		return
	}

	fmt.Printf("[CatchUp] Sending catch-up reminder for '%s' (%s-%s)\n",
		schedule.Activity, schedule.StartTime, schedule.EndTime)
	sendReminder(*schedule, now)
}

// parseChatID reads TELEGRAM_CHAT_ID from env.

func parseChatID() (int64, error) {
	val := os.Getenv("TELEGRAM_CHAT_ID")
	if val == "" {
		return 6616220735, nil
	}
	return strconv.ParseInt(val, 10, 64)
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
