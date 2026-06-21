// Package timeutil provides a collection of time-related helper functions for
// use throughout the Grafana backend. It supplements the standard library's
// time package with utilities for human-readable formatting, flexible parsing,
// and common time-range arithmetic.
package timeutil

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// Common timestamp layout constants used in Grafana's backend and API.
const (
	// LayoutISO is the canonical ISO-8601 / RFC-3339 layout used in JSON APIs.
	LayoutISO = time.RFC3339

	// LayoutDateOnly is a date-only layout (YYYY-MM-DD).
	LayoutDateOnly = "2006-01-02"

	// LayoutDateTime is a layout that includes date and time without a timezone
	// suffix (YYYY-MM-DD HH:MM:SS).
	LayoutDateTime = "2006-01-02 15:04:05"

	// LayoutDateTimeMS extends LayoutDateTime with millisecond precision.
	LayoutDateTimeMS = "2006-01-02 15:04:05.000"
)

// ParseAny tries to parse s against a series of well-known layouts and returns
// the first successful result in UTC. An error is returned only if none of the
// attempted layouts matches.
//
// Layouts tried (in order):
//   1. time.RFC3339Nano
//   2. time.RFC3339
//   3. LayoutDateTime
//   4. LayoutDateOnly
func ParseAny(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		LayoutDateTime,
		LayoutDateOnly,
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("timeutil: cannot parse %q with any known layout", s)
}

// MustParseAny is like ParseAny but panics on error. Useful for hard-coded
// sentinel values in tests and package-level variables.
func MustParseAny(s string) time.Time {
	t, err := ParseAny(s)
	if err != nil {
		panic(err)
	}
	return t
}

// FormatEpochMS formats a Unix timestamp expressed in milliseconds as an
// RFC-3339 string in UTC.
func FormatEpochMS(epochMs int64) string {
	return time.Unix(epochMs/1000, (epochMs%1000)*int64(time.Millisecond)).UTC().Format(LayoutISO)
}

// EpochMS returns the Unix timestamp of t in milliseconds.
func EpochMS(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// StartOfDay returns the instant at midnight (00:00:00.000) of the day that
// contains t, in the same location as t.
func StartOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the last nanosecond of the day that contains t, in the same
// location as t.
func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// StartOfWeek returns midnight on Monday of the ISO week that contains t, in
// the same location as t.
func StartOfWeek(t time.Time) time.Time {
	wd := t.Weekday()
	// Shift so that Monday == 0 … Sunday == 6.
	offset := int(wd) - int(time.Monday)
	if offset < 0 {
		offset += 7
	}
	return StartOfDay(t.AddDate(0, 0, -offset))
}

// StartOfMonth returns midnight on the first day of the month that contains t,
// in the same location as t.
func StartOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the last nanosecond of the last day of the month that
// contains t, in the same location as t.
func EndOfMonth(t time.Time) time.Time {
	return EndOfDay(StartOfMonth(t).AddDate(0, 1, -1))
}

// StartOfYear returns midnight on January 1st of the year that contains t, in
// the same location as t.
func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
}

// Age returns a human-readable description of how long ago t occurred relative
// to now. Results are intentionally coarse:
//
//   < 1 minute  → "just now"
//   < 1 hour    → "N minutes ago"
//   < 24 hours  → "N hours ago"
//   < 30 days   → "N days ago"
//   < 12 months → "N months ago"
//   otherwise   → "N years ago"
//
// If t is in the future the result is prefixed with "in ".
func Age(t time.Time) string {
	delta := time.Since(t)
	future := delta < 0
	if future {
		delta = -delta
	}

	var label string
	switch {
	case delta < time.Minute:
		label = "just now"
	case delta < time.Hour:
		mins := int(math.Round(delta.Minutes()))
		label = plural(mins, "minute") + " ago"
	case delta < 24*time.Hour:
		hrs := int(math.Round(delta.Hours()))
		label = plural(hrs, "hour") + " ago"
	case delta < 30*24*time.Hour:
		days := int(math.Round(delta.Hours() / 24))
		label = plural(days, "day") + " ago"
	case delta < 365*24*time.Hour:
		months := int(math.Round(delta.Hours() / 24 / 30))
		label = plural(months, "month") + " ago"
	default:
		years := int(math.Round(delta.Hours() / 24 / 365))
		label = plural(years, "year") + " ago"
	}

	if future && label != "just now" {
		label = "in " + strings.TrimSuffix(label, " ago")
	}
	return label
}

func plural(n int, unit string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", unit)
	}
	return fmt.Sprintf("%d %ss", n, unit)
}

// IsSameDay reports whether a and b fall on the same calendar date in UTC.
func IsSameDay(a, b time.Time) bool {
	ay, am, ad := a.UTC().Date()
	by, bm, bd := b.UTC().Date()
	return ay == by && am == bm && ad == bd
}

// Clamp returns t clamped to the interval [lo, hi]. If t is before lo it
// returns lo; if t is after hi it returns hi; otherwise it returns t unchanged.
// lo must be ≤ hi, or Clamp panics.
func Clamp(t, lo, hi time.Time) (time.Time, error) {
	if hi.Before(lo) {
		return time.Time{}, errors.New("timeutil: Clamp called with lo > hi")
	}
	if t.Before(lo) {
		return lo, nil
	}
	if t.After(hi) {
		return hi, nil
	}
	return t, nil
}

// BusinessDaysUntil counts the number of weekdays (Monday–Friday) between
// start (inclusive) and end (exclusive). A negative result means end is before
// start.
func BusinessDaysUntil(start, end time.Time) int {
	if start.Equal(end) {
		return 0
	}

	negative := end.Before(start)
	if negative {
		start, end = end, start
	}

	count := 0
	cur := StartOfDay(start)
	target := StartOfDay(end)
	for cur.Before(target) {
		wd := cur.Weekday()
		if wd != time.Saturday && wd != time.Sunday {
			count++
		}
		cur = cur.AddDate(0, 0, 1)
	}
	if negative {
		return -count
	}
	return count
}
