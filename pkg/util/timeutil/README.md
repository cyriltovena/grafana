# pkg/util/timeutil

`timeutil` is a Go package that augments the standard `time` library with
utilities commonly needed across the Grafana backend.

## Sub-packages

| File | Contents |
|------|----------|
| `timeutil.go` | Parsing, formatting, day/week/month boundaries, age strings, clamping |
| `duration.go` | Extended duration parsing (`d` / `w` units), human-readable formatting, rounding |
| `range.go` | `Range` type for closed time intervals with set-algebra operations |

---

## Quick-start

```go
import "github.com/grafana/grafana/pkg/util/timeutil"
```

### Parsing

```go
// Tries RFC3339Nano → RFC3339 → "2006-01-02 15:04:05" → "2006-01-02".
t, err := timeutil.ParseAny("2021-06-15 12:30:00")

// Panics on invalid input — useful for package-level variables.
sentinel := timeutil.MustParseAny("2021-01-01T00:00:00Z")
```

### Epoch milliseconds

```go
ms  := timeutil.EpochMS(time.Now())          // int64
str := timeutil.FormatEpochMS(1609459200000) // "2021-01-01T00:00:00Z"
```

### Boundary helpers

```go
timeutil.StartOfDay(t)    // 00:00:00.000 on the same day
timeutil.EndOfDay(t)      // 23:59:59.999999999 on the same day
timeutil.StartOfWeek(t)   // Monday midnight
timeutil.StartOfMonth(t)  // 1st of month, midnight
timeutil.EndOfMonth(t)    // last moment of the month
timeutil.StartOfYear(t)   // Jan 1 midnight
```

### Human-readable age

```go
timeutil.Age(t) // "just now", "3 hours ago", "2 days ago", …
```

### Clamping and business days

```go
clamped, err := timeutil.Clamp(t, lo, hi)
days := timeutil.BusinessDaysUntil(start, end)
```

---

### Extended duration parsing

```go
d, err := timeutil.ParseExtendedDuration("1w2d3h")
// = 7*24h + 2*24h + 3h = 228h

timeutil.FormatDurationHuman(228 * time.Hour) // "1w 2d" (two largest components)
timeutil.RoundDuration(90*time.Second, time.Minute) // 2m
```

---

### Range arithmetic

```go
r, _ := timeutil.NewRange(from, to)

r.Contains(t)            // bool
r.Overlaps(other)        // bool
r.Intersection(other)    // (Range, bool)
r.Union(other)           // Range
r.Shift(24*time.Hour)    // Range displaced by +24 h
r.Expand(time.Hour)      // Range widened by 1 h on each side
r.Split(7)               // []Range — 7 equal sub-ranges
r.Duration()             // time.Duration

// Convenience constructors
timeutil.RangeForDay(t)   // full calendar day
timeutil.RangeForMonth(t) // full calendar month
timeutil.LastN(time.Hour) // last 1 h ending now
```

---

## Running the tests

```sh
go test ./pkg/util/timeutil/... -v
```

All functions are unit-tested in the corresponding `_test.go` files.
