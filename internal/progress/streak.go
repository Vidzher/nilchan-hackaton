package progress

import "time"

func NextStreak(current int, previous *time.Time, completedAt time.Time) int {
	if previous == nil {
		return 1
	}

	currentDay := utcDay(completedAt)
	previousDay := utcDay(*previous)
	switch {
	case previousDay.Equal(currentDay):
		return current
	case previousDay.Equal(currentDay.AddDate(0, 0, -1)):
		return current + 1
	default:
		return 1
	}
}

func utcDay(value time.Time) time.Time {
	value = value.UTC()
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}
