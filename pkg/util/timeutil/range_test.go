package timeutil

import (
	"testing"
	"time"
)

func ref(s string) time.Time {
	t, err := time.Parse(LayoutISO, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestNewRange(t *testing.T) {
	from := ref("2021-01-01T00:00:00Z")
	to := ref("2021-12-31T23:59:59Z")

	r, err := NewRange(from, to)
	if err != nil {
		t.Fatalf("NewRange unexpected error: %v", err)
	}
	if !r.From.Equal(from) || !r.To.Equal(to) {
		t.Errorf("NewRange fields mismatch")
	}

	// Reversed should error.
	if _, err := NewRange(to, from); err == nil {
		t.Error("expected error for reversed range, got nil")
	}
}

func TestRangeDuration(t *testing.T) {
	from := ref("2021-01-01T00:00:00Z")
	to := ref("2021-01-02T00:00:00Z") // exactly 24 h later
	r := Range{From: from, To: to}
	if r.Duration() != 24*time.Hour {
		t.Errorf("Duration() = %v; want 24h", r.Duration())
	}
}

func TestRangeContains(t *testing.T) {
	r := Range{
		From: ref("2021-01-01T00:00:00Z"),
		To:   ref("2021-12-31T23:59:59Z"),
	}
	if !r.Contains(ref("2021-06-15T12:00:00Z")) {
		t.Error("Contains(mid) = false; want true")
	}
	if r.Contains(ref("2020-12-31T23:59:59Z")) {
		t.Error("Contains(before) = true; want false")
	}
	if r.Contains(ref("2022-01-01T00:00:00Z")) {
		t.Error("Contains(after) = true; want false")
	}
	// Boundary points are inclusive.
	if !r.Contains(r.From) || !r.Contains(r.To) {
		t.Error("boundary points should be contained")
	}
}

func TestRangeOverlaps(t *testing.T) {
	a := Range{From: ref("2021-01-01T00:00:00Z"), To: ref("2021-06-30T23:59:59Z")}
	b := Range{From: ref("2021-04-01T00:00:00Z"), To: ref("2021-12-31T23:59:59Z")}
	c := Range{From: ref("2022-01-01T00:00:00Z"), To: ref("2022-06-30T00:00:00Z")}

	if !a.Overlaps(b) {
		t.Error("a.Overlaps(b) = false; want true")
	}
	if a.Overlaps(c) {
		t.Error("a.Overlaps(c) = true; want false")
	}
	// Adjacent (sharing a boundary point) counts as overlapping.
	adj := Range{From: ref("2021-06-30T23:59:59Z"), To: ref("2021-12-31T00:00:00Z")}
	if !a.Overlaps(adj) {
		t.Error("a.Overlaps(adj) = false; want true (shared boundary)")
	}
}

func TestRangeIntersection(t *testing.T) {
	a := Range{From: ref("2021-01-01T00:00:00Z"), To: ref("2021-06-30T00:00:00Z")}
	b := Range{From: ref("2021-04-01T00:00:00Z"), To: ref("2021-12-31T00:00:00Z")}

	inter, ok := a.Intersection(b)
	if !ok {
		t.Fatal("Intersection(a, b) returned ok=false; want true")
	}
	if !inter.From.Equal(b.From) || !inter.To.Equal(a.To) {
		t.Errorf("Intersection = %v; want [%v, %v]", inter, b.From, a.To)
	}

	// Non-overlapping.
	c := Range{From: ref("2022-01-01T00:00:00Z"), To: ref("2022-06-30T00:00:00Z")}
	if _, ok := a.Intersection(c); ok {
		t.Error("Intersection(a, c) ok=true; want false")
	}
}

func TestRangeUnion(t *testing.T) {
	a := Range{From: ref("2021-01-01T00:00:00Z"), To: ref("2021-06-30T00:00:00Z")}
	b := Range{From: ref("2021-04-01T00:00:00Z"), To: ref("2021-12-31T00:00:00Z")}

	u := a.Union(b)
	if !u.From.Equal(a.From) || !u.To.Equal(b.To) {
		t.Errorf("Union = %v; want [%v, %v]", u, a.From, b.To)
	}
}

func TestRangeShift(t *testing.T) {
	r := Range{From: ref("2021-01-01T00:00:00Z"), To: ref("2021-01-02T00:00:00Z")}
	shifted := r.Shift(24 * time.Hour)
	want := Range{From: ref("2021-01-02T00:00:00Z"), To: ref("2021-01-03T00:00:00Z")}
	if !shifted.From.Equal(want.From) || !shifted.To.Equal(want.To) {
		t.Errorf("Shift = %v; want %v", shifted, want)
	}
}

func TestRangeExpand(t *testing.T) {
	r := Range{From: ref("2021-01-02T00:00:00Z"), To: ref("2021-01-03T00:00:00Z")}
	expanded := r.Expand(24 * time.Hour)
	if !expanded.From.Equal(ref("2021-01-01T00:00:00Z")) {
		t.Errorf("Expand From = %v; want 2021-01-01", expanded.From)
	}
	if !expanded.To.Equal(ref("2021-01-04T00:00:00Z")) {
		t.Errorf("Expand To = %v; want 2021-01-04", expanded.To)
	}
}

func TestRangeSplit(t *testing.T) {
	from := ref("2021-01-01T00:00:00Z")
	to := ref("2021-01-04T00:00:00Z") // 72 h total
	r := Range{From: from, To: to}

	parts, err := r.Split(3)
	if err != nil {
		t.Fatalf("Split(3) unexpected error: %v", err)
	}
	if len(parts) != 3 {
		t.Fatalf("Split(3) returned %d parts; want 3", len(parts))
	}
	// Each part should be 24 h.
	for i, p := range parts {
		if p.Duration() != 24*time.Hour {
			t.Errorf("part %d duration = %v; want 24h", i, p.Duration())
		}
	}
	// First and last boundaries should match the original.
	if !parts[0].From.Equal(from) {
		t.Errorf("parts[0].From = %v; want %v", parts[0].From, from)
	}
	if !parts[2].To.Equal(to) {
		t.Errorf("parts[2].To = %v; want %v", parts[2].To, to)
	}

	// Error on n <= 0.
	if _, err := r.Split(0); err == nil {
		t.Error("Split(0) expected error, got nil")
	}
}

func TestRangeForDay(t *testing.T) {
	ref := time.Date(2021, 6, 15, 14, 30, 0, 0, time.UTC)
	r := RangeForDay(ref)
	if r.From.Hour() != 0 {
		t.Errorf("RangeForDay From hour = %d; want 0", r.From.Hour())
	}
	if r.To.Hour() != 23 {
		t.Errorf("RangeForDay To hour = %d; want 23", r.To.Hour())
	}
}

func TestLastN(t *testing.T) {
	before := time.Now().UTC()
	r := LastN(time.Hour)
	after := time.Now().UTC()

	if r.Duration() != time.Hour {
		t.Errorf("LastN(1h) duration = %v; want 1h", r.Duration())
	}
	if r.To.Before(before) || r.To.After(after) {
		t.Errorf("LastN To is out of expected bounds")
	}
}
