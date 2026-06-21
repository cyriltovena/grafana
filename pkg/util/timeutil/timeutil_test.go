package timeutil

import (
	"testing"
	"time"
)

func TestParseAny(t *testing.T) {
	cases := []struct {
		input string
		want  string // expected UTC RFC3339
	}{
		{"2021-06-15T12:30:00Z", "2021-06-15T12:30:00Z"},
		{"2021-06-15T12:30:00+02:00", "2021-06-15T10:30:00Z"},
		{"2021-06-15T12:30:00.000000001Z", "2021-06-15T12:30:00Z"},
		{"2021-06-15 12:30:00", "2021-06-15T12:30:00Z"},
		{"2021-06-15", "2021-06-15T00:00:00Z"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseAny(tc.input)
			if err != nil {
				t.Fatalf("ParseAny(%q) returned unexpected error: %v", tc.input, err)
			}
			if got.Format(LayoutISO) != tc.want {
				t.Errorf("ParseAny(%q) = %q; want %q", tc.input, got.Format(LayoutISO), tc.want)
			}
		})
	}
}

func TestParseAnyError(t *testing.T) {
	if _, err := ParseAny("not-a-date"); err == nil {
		t.Error("expected error for invalid input, got nil")
	}
	if _, err := ParseAny(""); err == nil {
		t.Error("expected error for empty string, got nil")
	}
}

func TestMustParseAny(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected MustParseAny to panic on bad input")
		}
	}()
	MustParseAny("bad input")
}

func TestEpochMS(t *testing.T) {
	ref := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	ms := EpochMS(ref)
	back := time.Unix(ms/1000, (ms%1000)*int64(time.Millisecond)).UTC()
	if !back.Equal(ref) {
		t.Errorf("EpochMS round-trip failed: got %v, want %v", back, ref)
	}
}

func TestFormatEpochMS(t *testing.T) {
	ms := int64(1609459200000) // 2021-01-01T00:00:00Z
	got := FormatEpochMS(ms)
	if got != "2021-01-01T00:00:00Z" {
		t.Errorf("FormatEpochMS(%d) = %q; want %q", ms, got, "2021-01-01T00:00:00Z")
	}
}

func TestStartEndOfDay(t *testing.T) {
	loc := time.UTC
	ref := time.Date(2021, 6, 15, 14, 30, 0, 0, loc)

	start := StartOfDay(ref)
	if start.Hour() != 0 || start.Minute() != 0 || start.Second() != 0 {
		t.Errorf("StartOfDay unexpected: %v", start)
	}
	if start.Day() != 15 {
		t.Errorf("StartOfDay day mismatch: got %d", start.Day())
	}

	end := EndOfDay(ref)
	if end.Hour() != 23 || end.Minute() != 59 || end.Second() != 59 {
		t.Errorf("EndOfDay unexpected: %v", end)
	}
}

func TestStartOfWeek(t *testing.T) {
	// 2021-06-16 is a Wednesday; week should start on Monday 2021-06-14.
	wednesday := time.Date(2021, 6, 16, 10, 0, 0, 0, time.UTC)
	week := StartOfWeek(wednesday)
	if week.Weekday() != time.Monday {
		t.Errorf("StartOfWeek weekday = %v; want Monday", week.Weekday())
	}
	if week.Day() != 14 {
		t.Errorf("StartOfWeek day = %d; want 14", week.Day())
	}

	// Test on a Sunday.
	sunday := time.Date(2021, 6, 20, 0, 0, 0, 0, time.UTC)
	sw := StartOfWeek(sunday)
	if sw.Day() != 14 {
		t.Errorf("StartOfWeek(sunday) day = %d; want 14", sw.Day())
	}

	// Test on a Monday (should return same day).
	monday := time.Date(2021, 6, 14, 8, 0, 0, 0, time.UTC)
	sm := StartOfWeek(monday)
	if sm.Day() != 14 {
		t.Errorf("StartOfWeek(monday) day = %d; want 14", sm.Day())
	}
}

func TestStartEndOfMonth(t *testing.T) {
	ref := time.Date(2021, 3, 15, 0, 0, 0, 0, time.UTC)

	start := StartOfMonth(ref)
	if start.Day() != 1 || start.Month() != time.March {
		t.Errorf("StartOfMonth unexpected: %v", start)
	}

	end := EndOfMonth(ref)
	if end.Day() != 31 || end.Month() != time.March {
		t.Errorf("EndOfMonth unexpected: %v", end)
	}

	// February non-leap year.
	feb := time.Date(2021, 2, 10, 0, 0, 0, 0, time.UTC)
	if EndOfMonth(feb).Day() != 28 {
		t.Errorf("EndOfMonth Feb 2021 day = %d; want 28", EndOfMonth(feb).Day())
	}

	// February leap year.
	febLeap := time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC)
	if EndOfMonth(febLeap).Day() != 29 {
		t.Errorf("EndOfMonth Feb 2020 day = %d; want 29", EndOfMonth(febLeap).Day())
	}
}

func TestStartOfYear(t *testing.T) {
	ref := time.Date(2021, 6, 15, 12, 0, 0, 0, time.UTC)
	y := StartOfYear(ref)
	if y.Year() != 2021 || y.Month() != time.January || y.Day() != 1 {
		t.Errorf("StartOfYear unexpected: %v", y)
	}
}

func TestIsSameDay(t *testing.T) {
	a := time.Date(2021, 6, 15, 8, 0, 0, 0, time.UTC)
	b := time.Date(2021, 6, 15, 23, 59, 59, 0, time.UTC)
	c := time.Date(2021, 6, 16, 0, 0, 0, 0, time.UTC)

	if !IsSameDay(a, b) {
		t.Error("IsSameDay(a, b) = false; want true")
	}
	if IsSameDay(a, c) {
		t.Error("IsSameDay(a, c) = true; want false")
	}
}

func TestClamp(t *testing.T) {
	lo := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	hi := time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)

	// Within range.
	mid := time.Date(2021, 6, 15, 0, 0, 0, 0, time.UTC)
	got, err := Clamp(mid, lo, hi)
	if err != nil {
		t.Fatalf("Clamp returned unexpected error: %v", err)
	}
	if !got.Equal(mid) {
		t.Errorf("Clamp(mid) = %v; want %v", got, mid)
	}

	// Below lo.
	before := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)
	got, err = Clamp(before, lo, hi)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(lo) {
		t.Errorf("Clamp(before) = %v; want lo", got)
	}

	// Above hi.
	after := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	got, err = Clamp(after, lo, hi)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(hi) {
		t.Errorf("Clamp(after) = %v; want hi", got)
	}

	// Invalid range.
	if _, err := Clamp(mid, hi, lo); err == nil {
		t.Error("expected error for lo > hi, got nil")
	}
}

func TestBusinessDaysUntil(t *testing.T) {
	// Mon 2021-06-14 → Fri 2021-06-18 = 5 business days.
	mon := time.Date(2021, 6, 14, 0, 0, 0, 0, time.UTC)
	fri := time.Date(2021, 6, 18, 0, 0, 0, 0, time.UTC)
	if got := BusinessDaysUntil(mon, fri); got != 4 {
		// Mon→Fri exclusive means Mon,Tue,Wed,Thu counted = 4
		t.Errorf("BusinessDaysUntil(mon, fri) = %d; want 4", got)
	}

	// Same day.
	if got := BusinessDaysUntil(mon, mon); got != 0 {
		t.Errorf("BusinessDaysUntil(same) = %d; want 0", got)
	}

	// Negative (end before start).
	if got := BusinessDaysUntil(fri, mon); got >= 0 {
		t.Errorf("BusinessDaysUntil(fri, mon) = %d; want negative", got)
	}
}

func TestAge(t *testing.T) {
	now := time.Now()

	if Age(now.Add(-30*time.Second)) != "just now" {
		t.Errorf("expected 'just now' for 30 seconds ago")
	}
	if got := Age(now.Add(-2 * time.Minute)); got != "2 minutes ago" {
		t.Errorf("expected '2 minutes ago', got %q", got)
	}
	if got := Age(now.Add(-3 * time.Hour)); got != "3 hours ago" {
		t.Errorf("expected '3 hours ago', got %q", got)
	}
	if got := Age(now.Add(-14 * 24 * time.Hour)); got != "14 days ago" {
		t.Errorf("expected '14 days ago', got %q", got)
	}
}
