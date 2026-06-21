package timeutil

import (
	"errors"
	"fmt"
	"time"
)

// Range represents a closed time interval [From, To].
type Range struct {
	From time.Time
	To   time.Time
}

// ErrInvalidRange is returned when a Range has From after To.
var ErrInvalidRange = errors.New("timeutil: range From must not be after To")

// NewRange creates a Range and validates that From ≤ To.
func NewRange(from, to time.Time) (Range, error) {
	if from.After(to) {
		return Range{}, ErrInvalidRange
	}
	return Range{From: from, To: to}, nil
}

// Duration returns the duration covered by the range (To − From).
func (r Range) Duration() time.Duration {
	return r.To.Sub(r.From)
}

// Contains reports whether t falls within the range [From, To] (inclusive on
// both ends).
func (r Range) Contains(t time.Time) bool {
	return !t.Before(r.From) && !t.After(r.To)
}

// Overlaps reports whether r and other share at least one point in time.
func (r Range) Overlaps(other Range) bool {
	return !r.To.Before(other.From) && !other.To.Before(r.From)
}

// Intersection returns the overlap between r and other, and a boolean
// indicating whether an overlap exists.
func (r Range) Intersection(other Range) (Range, bool) {
	from := r.From
	if other.From.After(from) {
		from = other.From
	}
	to := r.To
	if other.To.Before(to) {
		to = other.To
	}
	if from.After(to) {
		return Range{}, false
	}
	return Range{From: from, To: to}, true
}

// Union returns the smallest Range that covers both r and other. The ranges do
// not need to overlap.
func (r Range) Union(other Range) Range {
	from := r.From
	if other.From.Before(from) {
		from = other.From
	}
	to := r.To
	if other.To.After(to) {
		to = other.To
	}
	return Range{From: from, To: to}
}

// Shift returns a new Range displaced by d relative to r.
func (r Range) Shift(d time.Duration) Range {
	return Range{From: r.From.Add(d), To: r.To.Add(d)}
}

// Expand returns a new Range that extends r by d on each side, i.e.
// [From − d, To + d].
func (r Range) Expand(d time.Duration) Range {
	return Range{From: r.From.Add(-d), To: r.To.Add(d)}
}

// Split divides r into n approximately equal sub-ranges. If n ≤ 0 an error is
// returned. The sub-ranges are contiguous and together cover exactly r.
func (r Range) Split(n int) ([]Range, error) {
	if n <= 0 {
		return nil, errors.New("timeutil: Split count must be > 0")
	}
	total := r.To.Sub(r.From)
	sliceNs := total.Nanoseconds() / int64(n)

	ranges := make([]Range, n)
	for i := 0; i < n; i++ {
		start := r.From.Add(time.Duration(int64(i) * sliceNs))
		var end time.Time
		if i == n-1 {
			end = r.To
		} else {
			end = r.From.Add(time.Duration(int64(i+1) * sliceNs))
		}
		ranges[i] = Range{From: start, To: end}
	}
	return ranges, nil
}

// String returns a human-readable representation of the range.
func (r Range) String() string {
	return fmt.Sprintf("[%s — %s]", r.From.UTC().Format(LayoutISO), r.To.UTC().Format(LayoutISO))
}

// RangeForDay returns the Range covering the full calendar day that contains t
// (midnight-to-midnight, inclusive) in the same location as t.
func RangeForDay(t time.Time) Range {
	return Range{From: StartOfDay(t), To: EndOfDay(t)}
}

// RangeForMonth returns the Range covering the full calendar month that
// contains t, in the same location as t.
func RangeForMonth(t time.Time) Range {
	return Range{From: StartOfMonth(t), To: EndOfMonth(t)}
}

// LastN returns a Range representing the last d of time ending at now (UTC).
func LastN(d time.Duration) Range {
	now := time.Now().UTC()
	return Range{From: now.Add(-d), To: now}
}
