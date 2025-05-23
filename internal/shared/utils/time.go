package utils

import (
	"time"
)

const (
	// Common time formats
	ISO8601     = "2006-01-02T15:04:05Z07:00"
	DateOnly    = "2006-01-02"
	TimeOnly    = "15:04:05"
	DateTimeStr = "2006-01-02 15:04:05"
)

// TimeZoneVN represents the Vietnam timezone.
var TimeZoneVN = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60)

// Now returns the current time in Vietnam timezone.
func Now() time.Time {
	return time.Now().In(TimeZoneVN)
}

// ParseTime parses a time string using the specified layout.
// Returns error if the parsing fails.
func ParseTime(value string, layout string) (time.Time, error) {
	return time.Parse(layout, value)
}

// FormatTime formats a time.Time using the specified layout.
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// StartOfDay returns the start of the day (00:00:00) for the given time.
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day (23:59:59.999999999) for the given time.
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// StartOfMonth returns the start of the month for the given time.
func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the end of the month for the given time.
func EndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month+1, 0, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// IsWeekend returns true if the given time is a weekend (Saturday or Sunday).
func IsWeekend(t time.Time) bool {
	day := t.Weekday()
	return day == time.Saturday || day == time.Sunday
}

// AddWorkDays adds the specified number of work days (excluding weekends) to the given time.
func AddWorkDays(t time.Time, days int) time.Time {
	// Handle negative days
	negative := days < 0
	if negative {
		days = -days
	}

	result := t
	for days > 0 {
		if negative {
			result = result.AddDate(0, 0, -1)
		} else {
			result = result.AddDate(0, 0, 1)
		}

		if !IsWeekend(result) {
			days--
		}
	}
	return result
}

// DurationToMilliseconds converts a time.Duration to milliseconds.
func DurationToMilliseconds(d time.Duration) int64 {
	return d.Milliseconds()
}

// MillisecondsToDuration converts milliseconds to time.Duration.
func MillisecondsToDuration(ms int64) time.Duration {
	return time.Duration(ms) * time.Millisecond
}
