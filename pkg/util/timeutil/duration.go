package timeutil

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// duration unit multipliers expressed as time.Duration values.
var durationUnits = []struct {
	suffix   string
	value    time.Duration
}{
	{"ns", time.Nanosecond},
	{"us", time.Microsecond},
	{"µs", time.Microsecond},
	{"ms", time.Millisecond},
	{"s", time.Second},
	{"m", time.Minute},
	{"h", time.Hour},
	{"d", 24 * time.Hour},
	{"w", 7 * 24 * time.Hour},
}

// durationPattern matches an optional leading sign, one or more groups of
// <number><unit>, e.g. "1h30m", "-2d12h", "90s".
var durationPattern = regexp.MustCompile(`^([+-]?)(\d+[a-zµ]+)+$`)

// ErrInvalidDuration is returned by ParseExtendedDuration when the input
// cannot be parsed.
var ErrInvalidDuration = errors.New("timeutil: invalid duration string")

// ParseExtendedDuration is a superset of time.ParseDuration that additionally
// recognises "d" (day = 24 h) and "w" (week = 7 d) suffixes.
//
// Examples:
//
//   "1w"       → 168h
//   "2d12h"    → 60h
//   "-3d"      → -72h
//   "1w2d3h4m" → 195h4m
//
// All standard time.ParseDuration suffixes (ns, us/µs, ms, s, m, h) are
// also accepted.
func ParseExtendedDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, ErrInvalidDuration
	}

	negative := false
	if s[0] == '-' {
		negative = true
		s = s[1:]
	} else if s[0] == '+' {
		s = s[1:]
	}

	if s == "" {
		return 0, ErrInvalidDuration
	}

	// First, try the stdlib so standard strings like "1h30m" still work
	// without going through our slower regex-based path when no d/w is present.
	if !strings.ContainsAny(s, "dw") {
		d, err := time.ParseDuration(s)
		if err != nil {
			return 0, ErrInvalidDuration
		}
		if negative {
			return -d, nil
		}
		return d, nil
	}

	total := time.Duration(0)
	remaining := s
	for remaining != "" {
		// Read numeric part.
		i := 0
		for i < len(remaining) && remaining[i] >= '0' && remaining[i] <= '9' {
			i++
		}
		if i == 0 {
			return 0, ErrInvalidDuration
		}
		n, err := strconv.ParseInt(remaining[:i], 10, 64)
		if err != nil {
			return 0, ErrInvalidDuration
		}
		remaining = remaining[i:]

		// Read unit part.
		j := 0
		for j < len(remaining) && (remaining[j] < '0' || remaining[j] > '9') {
			j++
		}
		if j == 0 {
			return 0, ErrInvalidDuration
		}
		unit := remaining[:j]
		remaining = remaining[j:]

		found := false
		for _, u := range durationUnits {
			if u.suffix == unit {
				total += time.Duration(n) * u.value
				found = true
				break
			}
		}
		if !found {
			return 0, fmt.Errorf("timeutil: unknown duration unit %q", unit)
		}
	}

	if negative {
		return -total, nil
	}
	return total, nil
}

// FormatDurationHuman converts d into a human-readable string using the two
// largest non-zero components.
//
// Examples:
//
//   0               → "0s"
//   90*time.Second  → "1m 30s"
//   25*time.Hour    → "1d 1h"
//   7*24*time.Hour  → "1w"
func FormatDurationHuman(d time.Duration) string {
	if d == 0 {
		return "0s"
	}
	negative := d < 0
	if negative {
		d = -d
	}

	type component struct {
		label string
		value int64
	}

	weeks := int64(d / (7 * 24 * time.Hour))
	d -= time.Duration(weeks) * 7 * 24 * time.Hour

	days := int64(d / (24 * time.Hour))
	d -= time.Duration(days) * 24 * time.Hour

	hours := int64(d / time.Hour)
	d -= time.Duration(hours) * time.Hour

	minutes := int64(d / time.Minute)
	d -= time.Duration(minutes) * time.Minute

	seconds := int64(d / time.Second)

	components := []component{
		{"w", weeks},
		{"d", days},
		{"h", hours},
		{"m", minutes},
		{"s", seconds},
	}

	var parts []string
	for _, c := range components {
		if c.value > 0 {
			parts = append(parts, fmt.Sprintf("%d%s", c.value, c.label))
		}
		if len(parts) == 2 {
			break
		}
	}

	result := strings.Join(parts, " ")
	if result == "" {
		result = "0s"
	}
	if negative {
		result = "-" + result
	}
	return result
}

// RoundDuration rounds d to the nearest multiple of unit. It behaves like
// time.Round but accepts any time.Duration as the rounding unit.
//
// Example:
//
//   RoundDuration(90*time.Second, time.Minute) → 2*time.Minute
func RoundDuration(d, unit time.Duration) time.Duration {
	if unit <= 0 {
		return d
	}
	return (d + unit/2) / unit * unit
}

// TruncateDuration truncates d to a multiple of unit (rounding towards zero).
func TruncateDuration(d, unit time.Duration) time.Duration {
	if unit <= 0 {
		return d
	}
	return d / unit * unit
}
