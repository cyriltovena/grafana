// Package util provides miscellaneous helper utilities used throughout Grafana.
package util

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ---------------------------------------------------------------------------
// String helpers
// ---------------------------------------------------------------------------

// Truncate returns at most maxLen runes of s.  If the string is longer it is
// cut and a "…" suffix is appended (the suffix counts toward maxLen).
func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen == 1 {
		return "…"
	}
	return string(runes[:maxLen-1]) + "…"
}

// Capitalize returns s with its first Unicode letter upper-cased and the rest
// lower-cased.
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// ToSnakeCase converts a camelCase or PascalCase string to snake_case.
// e.g. "MyHTTPServer" → "my_http_server".
func ToSnakeCase(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				// Add underscore before an upper-case letter that follows a
				// lower-case letter, or that is the start of a new "word"
				// within an acronym sequence (upper followed by lower).
				prev := runes[i-1]
				if unicode.IsLower(prev) || unicode.IsDigit(prev) {
					b.WriteRune('_')
				} else if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					b.WriteRune('_')
				}
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ToKebabCase converts a camelCase, PascalCase, or snake_case string to
// kebab-case.  e.g. "myVariable_name" → "my-variable-name".
func ToKebabCase(s string) string {
	return strings.ReplaceAll(ToSnakeCase(s), "_", "-")
}

// ContainsIgnoreCase reports whether substr appears within s, ignoring ASCII
// case differences.
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// HasPrefixIgnoreCase reports whether s begins with prefix, ignoring ASCII
// case differences.
func HasPrefixIgnoreCase(s, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
}

// HasSuffixIgnoreCase reports whether s ends with suffix, ignoring ASCII case
// differences.
func HasSuffixIgnoreCase(s, suffix string) bool {
	return strings.HasSuffix(strings.ToLower(s), strings.ToLower(suffix))
}

// Repeat returns s repeated n times.  Returns an empty string when n ≤ 0.
func Repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(s, n)
}

// ReverseString returns the UTF-8 reversal of s.
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome reports whether s reads the same forwards and backwards
// (Unicode-aware, case-insensitive).
func IsPalindrome(s string) bool {
	lower := []rune(strings.ToLower(s))
	for i, j := 0, len(lower)-1; i < j; i, j = i+1, j-1 {
		if lower[i] != lower[j] {
			return false
		}
	}
	return true
}

// WordCount returns the number of whitespace-separated words in s.
func WordCount(s string) int {
	return len(strings.Fields(s))
}

// UniqueStrings returns a new slice containing only the first occurrence of
// each element in ss (order is preserved).
func UniqueStrings(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

// ContainsString reports whether target is present in the slice ss.
func ContainsString(ss []string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}

// FilterStrings returns a new slice of strings from ss for which predicate
// returns true.
func FilterStrings(ss []string, predicate func(string) bool) []string {
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		if predicate(s) {
			out = append(out, s)
		}
	}
	return out
}

// MapStrings applies f to every element of ss and returns the resulting slice.
func MapStrings(ss []string, f func(string) string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = f(s)
	}
	return out
}

// JoinNonEmpty joins only the non-empty elements of ss with sep.
func JoinNonEmpty(ss []string, sep string) string {
	filtered := FilterStrings(ss, func(s string) bool { return s != "" })
	return strings.Join(filtered, sep)
}

// PadLeft pads s on the left with padChar until its rune length is at least
// totalLen.
func PadLeft(s string, totalLen int, padChar rune) string {
	runes := []rune(s)
	for len(runes) < totalLen {
		runes = append([]rune{padChar}, runes...)
	}
	return string(runes)
}

// PadRight pads s on the right with padChar until its rune length is at least
// totalLen.
func PadRight(s string, totalLen int, padChar rune) string {
	runes := []rune(s)
	for len(runes) < totalLen {
		runes = append(runes, padChar)
	}
	return string(runes)
}

// ---------------------------------------------------------------------------
// Numeric / conversion helpers
// ---------------------------------------------------------------------------

// SafeParseInt parses s as a base-10 integer. On error it returns defaultVal.
func SafeParseInt(s string, defaultVal int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

// SafeParseFloat64 parses s as a float64. On error it returns defaultVal.
func SafeParseFloat64(s string, defaultVal float64) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal
	}
	return v
}

// SafeParseBool parses s as a boolean using strconv.ParseBool. On error it
// returns defaultVal.
func SafeParseBool(s string, defaultVal bool) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return defaultVal
	}
	return v
}

// ClampInt returns v clamped to the range [min, max].
func ClampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ClampFloat64 returns v clamped to the range [min, max].
func ClampFloat64(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// MaxInt returns the larger of x and y.
func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// AbsInt returns the absolute value of x.
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// IntToHex returns the lower-case hexadecimal string representation of n.
func IntToHex(n int) string {
	return fmt.Sprintf("%x", n)
}

// HexToInt parses a hexadecimal string (with or without leading "0x") and
// returns the integer value.  Returns 0 and a non-nil error on failure.
func HexToInt(s string) (int64, error) {
	s = strings.TrimPrefix(strings.ToLower(s), "0x")
	return strconv.ParseInt(s, 16, 64)
}

// ---------------------------------------------------------------------------
// Validation helpers
// ---------------------------------------------------------------------------

// IsAlphanumeric reports whether every rune in s is a letter or digit.
func IsAlphanumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsNumeric reports whether every rune in s is a decimal digit.
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsEmail performs a rudimentary e-mail address check: the string must
// contain exactly one "@" with non-empty local and domain parts, and the
// domain part must contain at least one ".".
func IsEmail(s string) bool {
	parts := strings.SplitN(s, "@", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	return strings.Contains(parts[1], ".")
}

// IsURL performs a rudimentary URL check: the string must start with "http://"
// or "https://" and have a non-empty host.
func IsURL(s string) bool {
	lower := strings.ToLower(s)
	var rest string
	switch {
	case strings.HasPrefix(lower, "https://"):
		rest = s[len("https://"):]
	case strings.HasPrefix(lower, "http://"):
		rest = s[len("http://"):]
	default:
		return false
	}
	return len(strings.TrimSpace(rest)) > 0
}
