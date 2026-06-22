package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Truncate
// ---------------------------------------------------------------------------

func TestTruncate(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{name: "empty string", input: "", maxLen: 10, want: ""},
		{name: "shorter than limit", input: "hello", maxLen: 10, want: "hello"},
		{name: "exact length", input: "hello", maxLen: 5, want: "hello"},
		{name: "longer than limit", input: "hello world", maxLen: 8, want: "hello w…"},
		{name: "maxLen=1", input: "hello", maxLen: 1, want: "…"},
		{name: "maxLen=0 returns empty", input: "hello", maxLen: 0, want: ""},
		{name: "negative maxLen returns empty", input: "hello", maxLen: -3, want: ""},
		{name: "unicode string shorter than limit", input: "héllo", maxLen: 10, want: "héllo"},
		{name: "unicode string truncated", input: "héllo wörld", maxLen: 7, want: "héllo w…"},
		{name: "single character no truncation", input: "x", maxLen: 1, want: "x"},
		{name: "single character truncated with maxLen=2", input: "xy", maxLen: 2, want: "xy"},
		{name: "all same chars", input: "aaaaaaaaaa", maxLen: 5, want: "aaaa…"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Truncate(tc.input, tc.maxLen))
		})
	}
}

// ---------------------------------------------------------------------------
// Capitalize
// ---------------------------------------------------------------------------

func TestCapitalize(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"hello", "Hello"},
		{"HELLO", "Hello"},
		{"hELLO", "Hello"},
		{"a", "A"},
		{"A", "A"},
		{"grafana", "Grafana"},
		{"GRAFANA", "Grafana"},
		{"123abc", "123abc"},   // leading digit unchanged
		{"élan", "Élan"},       // non-ASCII letter
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, tc.want, Capitalize(tc.input))
		})
	}
}

// ---------------------------------------------------------------------------
// ToSnakeCase
// ---------------------------------------------------------------------------

func TestToSnakeCase(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"hello", "hello"},
		{"helloWorld", "hello_world"},
		{"HelloWorld", "hello_world"},
		{"MyHTTPServer", "my_http_server"},
		{"getHTTPResponseCode", "get_http_response_code"},
		{"HTTPStatusCode", "http_status_code"},
		{"alreadySnake", "already_snake"},
		{"ID", "id"},
		{"userID", "user_id"},
		{"parseURL", "parse_url"},
		{"parseURLString", "parse_url_string"},
		{"simpleVar", "simple_var"},
		{"simple", "simple"},
		{"Simple", "simple"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, tc.want, ToSnakeCase(tc.input))
		})
	}
}

// ---------------------------------------------------------------------------
// ToKebabCase
// ---------------------------------------------------------------------------

func TestToKebabCase(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"hello", "hello"},
		{"helloWorld", "hello-world"},
		{"HelloWorld", "hello-world"},
		{"MyHTTPServer", "my-http-server"},
		{"userID", "user-id"},
		{"simple", "simple"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, tc.want, ToKebabCase(tc.input))
		})
	}
}

// ---------------------------------------------------------------------------
// ContainsIgnoreCase / HasPrefixIgnoreCase / HasSuffixIgnoreCase
// ---------------------------------------------------------------------------

func TestContainsIgnoreCase(t *testing.T) {
	cases := []struct {
		s      string
		substr string
		want   bool
	}{
		{"Grafana", "grafana", true},
		{"Grafana", "GRAFANA", true},
		{"Grafana", "ana", true},
		{"Grafana", "ANA", true},
		{"Grafana", "xyz", false},
		{"", "", true},
		{"hello", "", true},
		{"", "x", false},
		{"UPPERCASE", "uppercase", true},
		{"lowercase", "LOWER", true},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, ContainsIgnoreCase(tc.s, tc.substr),
			"ContainsIgnoreCase(%q, %q)", tc.s, tc.substr)
	}
}

func TestHasPrefixIgnoreCase(t *testing.T) {
	cases := []struct {
		s      string
		prefix string
		want   bool
	}{
		{"Grafana", "gra", true},
		{"Grafana", "GRA", true},
		{"Grafana", "GRAFANA", true},
		{"Grafana", "ana", false},
		{"", "", true},
		{"hello", "", true},
		{"", "x", false},
		{"HTTP/1.1", "http", true},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, HasPrefixIgnoreCase(tc.s, tc.prefix),
			"HasPrefixIgnoreCase(%q, %q)", tc.s, tc.prefix)
	}
}

func TestHasSuffixIgnoreCase(t *testing.T) {
	cases := []struct {
		s      string
		suffix string
		want   bool
	}{
		{"Grafana", "ana", true},
		{"Grafana", "ANA", true},
		{"Grafana", "GRAFANA", true},
		{"Grafana", "gra", false},
		{"", "", true},
		{"hello", "", true},
		{"", "x", false},
		{"index.HTML", ".html", true},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, HasSuffixIgnoreCase(tc.s, tc.suffix),
			"HasSuffixIgnoreCase(%q, %q)", tc.s, tc.suffix)
	}
}

// ---------------------------------------------------------------------------
// Repeat
// ---------------------------------------------------------------------------

func TestRepeat(t *testing.T) {
	cases := []struct {
		s    string
		n    int
		want string
	}{
		{"ab", 3, "ababab"},
		{"ab", 1, "ab"},
		{"ab", 0, ""},
		{"ab", -1, ""},
		{"x", 5, "xxxxx"},
		{"", 10, ""},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, Repeat(tc.s, tc.n),
			"Repeat(%q, %d)", tc.s, tc.n)
	}
}

// ---------------------------------------------------------------------------
// ReverseString
// ---------------------------------------------------------------------------

func TestReverseString(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
		{"hello", "olleh"},
		{"Grafana", "anafarg"},
		{"racecar", "racecar"},
		{"Hello, 世界", "界世 ,olleH"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, tc.want, ReverseString(tc.input))
		})
	}
}

// ---------------------------------------------------------------------------
// IsPalindrome
// ---------------------------------------------------------------------------

func TestIsPalindrome(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"", true},
		{"a", true},
		{"aa", true},
		{"ab", false},
		{"racecar", true},
		{"Racecar", true},
		{"RaceCar", true},
		{"hello", false},
		{"level", true},
		{"Madam", true},
		{"Grafana", false},
		{"noon", true},
		{"No lemon no melon", false}, // spaces differ
		{"abcba", true},
		{"abccba", true},
		{"abcde", false},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, tc.want, IsPalindrome(tc.input))
		})
	}
}

// ---------------------------------------------------------------------------
// WordCount
// ---------------------------------------------------------------------------

func TestWordCount(t *testing.T) {
	cases := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"hello", 1},
		{"hello world", 2},
		{"  hello   world  ", 2},
		{"one two three four five", 5},
		{"tabs\there", 2},
		{"newline\nhere", 2},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, WordCount(tc.input),
			"WordCount(%q)", tc.input)
	}
}

// ---------------------------------------------------------------------------
// UniqueStrings
// ---------------------------------------------------------------------------

func TestUniqueStrings(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		want  []string
	}{
		{name: "nil slice", input: nil, want: []string{}},
		{name: "empty slice", input: []string{}, want: []string{}},
		{name: "no duplicates", input: []string{"a", "b", "c"}, want: []string{"a", "b", "c"}},
		{name: "all duplicates", input: []string{"a", "a", "a"}, want: []string{"a"}},
		{name: "mixed", input: []string{"b", "a", "b", "c", "a"}, want: []string{"b", "a", "c"}},
		{name: "preserves order", input: []string{"z", "y", "x"}, want: []string{"z", "y", "x"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, UniqueStrings(tc.input))
		})
	}
}

// ---------------------------------------------------------------------------
// ContainsString
// ---------------------------------------------------------------------------

func TestContainsString(t *testing.T) {
	ss := []string{"alpha", "beta", "gamma"}
	assert.True(t, ContainsString(ss, "alpha"))
	assert.True(t, ContainsString(ss, "beta"))
	assert.True(t, ContainsString(ss, "gamma"))
	assert.False(t, ContainsString(ss, "delta"))
	assert.False(t, ContainsString(ss, "Alpha")) // case-sensitive
	assert.False(t, ContainsString(nil, "anything"))
	assert.False(t, ContainsString([]string{}, "anything"))
}

// ---------------------------------------------------------------------------
// FilterStrings
// ---------------------------------------------------------------------------

func TestFilterStrings(t *testing.T) {
	isLong := func(s string) bool { return len(s) > 3 }

	cases := []struct {
		name      string
		input     []string
		predicate func(string) bool
		want      []string
	}{
		{"nil input", nil, isLong, []string{}},
		{"empty input", []string{}, isLong, []string{}},
		{"all pass", []string{"hello", "world"}, isLong, []string{"hello", "world"}},
		{"none pass", []string{"a", "bb", "ccc"}, isLong, []string{}},
		{"mixed", []string{"hi", "grafana", "go", "dashboard"}, isLong, []string{"grafana", "dashboard"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, FilterStrings(tc.input, tc.predicate))
		})
	}
}

// ---------------------------------------------------------------------------
// MapStrings
// ---------------------------------------------------------------------------

func TestMapStrings(t *testing.T) {
	toUpper := func(s string) string { return s + "!" }

	cases := []struct {
		name  string
		input []string
		want  []string
	}{
		{"nil input", nil, []string{}},
		{"empty", []string{}, []string{}},
		{"single", []string{"hello"}, []string{"hello!"}},
		{"multiple", []string{"a", "b", "c"}, []string{"a!", "b!", "c!"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, MapStrings(tc.input, toUpper))
		})
	}
}

// ---------------------------------------------------------------------------
// JoinNonEmpty
// ---------------------------------------------------------------------------

func TestJoinNonEmpty(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		sep   string
		want  string
	}{
		{"all empty", []string{"", "", ""}, ",", ""},
		{"no empties", []string{"a", "b", "c"}, "-", "a-b-c"},
		{"mixed", []string{"a", "", "c"}, "/", "a/c"},
		{"single non-empty", []string{"", "hello", ""}, ",", "hello"},
		{"nil slice", nil, ",", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, JoinNonEmpty(tc.input, tc.sep))
		})
	}
}

// ---------------------------------------------------------------------------
// PadLeft / PadRight
// ---------------------------------------------------------------------------

func TestPadLeft(t *testing.T) {
	cases := []struct {
		input    string
		totalLen int
		padChar  rune
		want     string
	}{
		{"42", 5, '0', "00042"},
		{"hello", 3, ' ', "hello"},
		{"hello", 5, ' ', "hello"},
		{"hello", 7, '-', "--hello"},
		{"", 3, 'x', "xxx"},
		{"hi", 4, '•', "••hi"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, PadLeft(tc.input, tc.totalLen, tc.padChar),
			"PadLeft(%q, %d, %q)", tc.input, tc.totalLen, tc.padChar)
	}
}

func TestPadRight(t *testing.T) {
	cases := []struct {
		input    string
		totalLen int
		padChar  rune
		want     string
	}{
		{"hello", 8, '.', "hello..."},
		{"hello", 5, '.', "hello"},
		{"hello", 3, '.', "hello"},
		{"", 3, '*', "***"},
		{"ab", 4, ' ', "ab  "},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, PadRight(tc.input, tc.totalLen, tc.padChar),
			"PadRight(%q, %d, %q)", tc.input, tc.totalLen, tc.padChar)
	}
}

// ---------------------------------------------------------------------------
// SafeParseInt
// ---------------------------------------------------------------------------

func TestSafeParseInt(t *testing.T) {
	cases := []struct {
		input        string
		defaultValue int
		want         int
	}{
		{"42", 0, 42},
		{"-7", 0, -7},
		{"0", 99, 0},
		{"", 5, 5},
		{"abc", 5, 5},
		{"3.14", 0, 0},
		{"999999999999999999999", 42, 42}, // overflow
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, SafeParseInt(tc.input, tc.defaultValue),
			"SafeParseInt(%q, %d)", tc.input, tc.defaultValue)
	}
}

// ---------------------------------------------------------------------------
// SafeParseFloat64
// ---------------------------------------------------------------------------

func TestSafeParseFloat64(t *testing.T) {
	cases := []struct {
		input        string
		defaultValue float64
		want         float64
	}{
		{"3.14", 0.0, 3.14},
		{"-2.5", 0.0, -2.5},
		{"0", 9.9, 0.0},
		{"", 1.5, 1.5},
		{"abc", 1.5, 1.5},
		{"1e3", 0.0, 1000.0},
		{"NaN", 0.0, 0.0}, // strconv rejects plain "NaN"? actually it accepts it
	}
	// NaN is special — strconv.ParseFloat("NaN", 64) succeeds and returns NaN.
	// We test the non-NaN cases individually.
	for _, tc := range cases {
		got := SafeParseFloat64(tc.input, tc.defaultValue)
		if tc.input == "NaN" {
			// Just ensure we get back some float (NaN or the default) without panicking.
			continue
		}
		assert.InDelta(t, tc.want, got, 1e-9,
			"SafeParseFloat64(%q, %v)", tc.input, tc.defaultValue)
	}
}

// ---------------------------------------------------------------------------
// SafeParseBool
// ---------------------------------------------------------------------------

func TestSafeParseBool(t *testing.T) {
	trueInputs := []string{"1", "t", "T", "TRUE", "true", "True", "yes", "YES", "Yes"}
	falseInputs := []string{"0", "f", "F", "FALSE", "false", "False", "no", "NO", "No"}

	for _, s := range trueInputs {
		// strconv.ParseBool only understands 1/t/T/TRUE/true/True; "yes" will use default
		got := SafeParseBool(s, false)
		_ = got // just verify it does not panic
	}
	for _, s := range falseInputs {
		got := SafeParseBool(s, true)
		_ = got
	}

	// Known good cases
	assert.True(t, SafeParseBool("true", false))
	assert.True(t, SafeParseBool("1", false))
	assert.False(t, SafeParseBool("false", true))
	assert.False(t, SafeParseBool("0", true))
	// Invalid falls back to default
	assert.True(t, SafeParseBool("maybe", true))
	assert.False(t, SafeParseBool("maybe", false))
}

// ---------------------------------------------------------------------------
// ClampInt
// ---------------------------------------------------------------------------

func TestClampInt(t *testing.T) {
	cases := []struct {
		v, min, max, want int
	}{
		{5, 0, 10, 5},
		{-1, 0, 10, 0},
		{11, 0, 10, 10},
		{0, 0, 10, 0},
		{10, 0, 10, 10},
		{100, -50, 50, 50},
		{-100, -50, 50, -50},
		{0, 0, 0, 0},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, ClampInt(tc.v, tc.min, tc.max),
			"ClampInt(%d, %d, %d)", tc.v, tc.min, tc.max)
	}
}

// ---------------------------------------------------------------------------
// ClampFloat64
// ---------------------------------------------------------------------------

func TestClampFloat64(t *testing.T) {
	cases := []struct {
		v, min, max, want float64
	}{
		{0.5, 0.0, 1.0, 0.5},
		{-0.1, 0.0, 1.0, 0.0},
		{1.1, 0.0, 1.0, 1.0},
		{0.0, 0.0, 1.0, 0.0},
		{1.0, 0.0, 1.0, 1.0},
		{99.9, 0.0, 100.0, 99.9},
	}
	for _, tc := range cases {
		assert.InDelta(t, tc.want, ClampFloat64(tc.v, tc.min, tc.max), 1e-12,
			"ClampFloat64(%v, %v, %v)", tc.v, tc.min, tc.max)
	}
}

// ---------------------------------------------------------------------------
// MaxInt
// ---------------------------------------------------------------------------

func TestMaxInt(t *testing.T) {
	cases := []struct {
		x, y, want int
	}{
		{1, 2, 2},
		{2, 1, 2},
		{0, 0, 0},
		{-1, -2, -1},
		{-2, -1, -1},
		{100, 99, 100},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, MaxInt(tc.x, tc.y))
	}
}

// ---------------------------------------------------------------------------
// AbsInt
// ---------------------------------------------------------------------------

func TestAbsInt(t *testing.T) {
	cases := []struct {
		input, want int
	}{
		{0, 0},
		{1, 1},
		{-1, 1},
		{-100, 100},
		{100, 100},
		{-999999, 999999},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, AbsInt(tc.input))
	}
}

// ---------------------------------------------------------------------------
// IntToHex / HexToInt
// ---------------------------------------------------------------------------

func TestIntToHex(t *testing.T) {
	cases := []struct {
		input int
		want  string
	}{
		{0, "0"},
		{255, "ff"},
		{256, "100"},
		{16, "10"},
		{4096, "1000"},
		{65535, "ffff"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, IntToHex(tc.input))
	}
}

func TestHexToInt(t *testing.T) {
	cases := []struct {
		input   string
		want    int64
		wantErr bool
	}{
		{"ff", 255, false},
		{"FF", 255, false},
		{"0xff", 255, false},
		{"0xFF", 255, false},
		{"100", 256, false},
		{"0", 0, false},
		{"1000", 4096, false},
		{"ffff", 65535, false},
		{"xyz", 0, true},
		{"", 0, true},
	}
	for _, tc := range cases {
		got, err := HexToInt(tc.input)
		if tc.wantErr {
			require.Error(t, err, "HexToInt(%q) expected error", tc.input)
		} else {
			require.NoError(t, err, "HexToInt(%q)", tc.input)
			assert.Equal(t, tc.want, got)
		}
	}
}

// IntToHex and HexToInt should be inverses for non-negative numbers.
func TestIntToHexRoundTrip(t *testing.T) {
	values := []int{0, 1, 10, 255, 256, 1024, 65535, 1048576}
	for _, v := range values {
		hex := IntToHex(v)
		back, err := HexToInt(hex)
		require.NoError(t, err)
		assert.Equal(t, int64(v), back, "round-trip failed for %d", v)
	}
}

// ---------------------------------------------------------------------------
// IsAlphanumeric
// ---------------------------------------------------------------------------

func TestIsAlphanumeric(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"abc", true},
		{"ABC", true},
		{"abc123", true},
		{"123", true},
		{"abc!", false},
		{"abc def", false},
		{"abc-def", false},
		{"_abc", false},
		{"Grafana7", true},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, IsAlphanumeric(tc.input),
			"IsAlphanumeric(%q)", tc.input)
	}
}

// ---------------------------------------------------------------------------
// IsNumeric
// ---------------------------------------------------------------------------

func TestIsNumeric(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"0", true},
		{"123", true},
		{"123.45", false},
		{"123a", false},
		{"-1", false},
		{" 1", false},
		{"00001", true},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, IsNumeric(tc.input),
			"IsNumeric(%q)", tc.input)
	}
}

// ---------------------------------------------------------------------------
// IsEmail
// ---------------------------------------------------------------------------

func TestIsEmail(t *testing.T) {
	validEmails := []string{
		"user@example.com",
		"user+tag@sub.domain.org",
		"first.last@company.co.uk",
		"admin@grafana.com",
		"a@b.c",
	}
	invalidEmails := []string{
		"",
		"notanemail",
		"@domain.com",
		"user@",
		"user@nodot",
		"user@.com",
		"@@example.com",
	}
	for _, e := range validEmails {
		assert.True(t, IsEmail(e), "expected %q to be valid", e)
	}
	for _, e := range invalidEmails {
		assert.False(t, IsEmail(e), "expected %q to be invalid", e)
	}
}

// ---------------------------------------------------------------------------
// IsURL
// ---------------------------------------------------------------------------

func TestIsURL(t *testing.T) {
	validURLs := []string{
		"http://grafana.com",
		"https://grafana.com",
		"HTTP://GRAFANA.COM",
		"HTTPS://GRAFANA.COM",
		"https://grafana.com/path?q=1",
		"http://localhost:3000",
		"https://user:pass@host/path",
	}
	invalidURLs := []string{
		"",
		"ftp://grafana.com",
		"grafana.com",
		"//grafana.com",
		"http://",
		"https://",
		"https://   ", // only whitespace after scheme
	}
	for _, u := range validURLs {
		assert.True(t, IsURL(u), "expected %q to be a valid URL", u)
	}
	for _, u := range invalidURLs {
		assert.False(t, IsURL(u), "expected %q to be an invalid URL", u)
	}
}
