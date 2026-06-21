package timeutil

import (
	"testing"
	"time"
)

func TestParseExtendedDuration(t *testing.T) {
	cases := []struct {
		input string
		want  time.Duration
		isErr bool
	}{
		{"1w", 7 * 24 * time.Hour, false},
		{"2d", 2 * 24 * time.Hour, false},
		{"1d12h", 36 * time.Hour, false},
		{"1w2d3h", 7*24*time.Hour + 2*24*time.Hour + 3*time.Hour, false},
		{"-1d", -24 * time.Hour, false},
		{"+2d", 2 * 24 * time.Hour, false},
		{"30s", 30 * time.Second, false},
		{"1h30m", 90 * time.Minute, false},
		{"500ms", 500 * time.Millisecond, false},
		{"", 0, true},
		{"1x", 0, true},
		{"-", 0, true},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseExtendedDuration(tc.input)
			if tc.isErr {
				if err == nil {
					t.Errorf("ParseExtendedDuration(%q): expected error, got %v", tc.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseExtendedDuration(%q) unexpected error: %v", tc.input, err)
			}
			if got != tc.want {
				t.Errorf("ParseExtendedDuration(%q) = %v; want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestFormatDurationHuman(t *testing.T) {
	cases := []struct {
		d    time.Duration
		want string
	}{
		{0, "0s"},
		{time.Second, "1s"},
		{90 * time.Second, "1m 30s"},
		{time.Hour, "1h"},
		{25 * time.Hour, "1d 1h"},
		{7 * 24 * time.Hour, "1w"},
		{8 * 24 * time.Hour, "1w 1d"},
		{-(90 * time.Second), "-1m 30s"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.want, func(t *testing.T) {
			got := FormatDurationHuman(tc.d)
			if got != tc.want {
				t.Errorf("FormatDurationHuman(%v) = %q; want %q", tc.d, got, tc.want)
			}
		})
	}
}

func TestRoundDuration(t *testing.T) {
	if got := RoundDuration(90*time.Second, time.Minute); got != 2*time.Minute {
		t.Errorf("RoundDuration(90s, 1m) = %v; want 2m", got)
	}
	if got := RoundDuration(29*time.Second, time.Minute); got != 0 {
		t.Errorf("RoundDuration(29s, 1m) = %v; want 0", got)
	}
	// Zero unit should return d unchanged.
	if got := RoundDuration(90*time.Second, 0); got != 90*time.Second {
		t.Errorf("RoundDuration(90s, 0) = %v; want 90s", got)
	}
}

func TestTruncateDuration(t *testing.T) {
	if got := TruncateDuration(90*time.Second, time.Minute); got != time.Minute {
		t.Errorf("TruncateDuration(90s, 1m) = %v; want 1m", got)
	}
}
