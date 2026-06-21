package util

// ContainsString reports whether slice contains the given string.
func ContainsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// ContainsInt reports whether slice contains the given int.
func ContainsInt(slice []int, n int) bool {
	for _, v := range slice {
		if v == n {
			return true
		}
	}
	return false
}

// UniqueStrings returns a new slice containing the elements of input with
// duplicates removed. Order is preserved (first occurrence wins).
func UniqueStrings(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	out := make([]string, 0, len(input))
	for _, s := range input {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

// FilterStrings returns a new slice containing only the elements of input for
// which keep returns true.
func FilterStrings(input []string, keep func(string) bool) []string {
	out := make([]string, 0, len(input))
	for _, s := range input {
		if keep(s) {
			out = append(out, s)
		}
	}
	return out
}

// MapStrings applies f to every element of input and returns the results.
func MapStrings(input []string, f func(string) string) []string {
	out := make([]string, len(input))
	for i, s := range input {
		out[i] = f(s)
	}
	return out
}

// Chunk splits input into consecutive sub-slices of at most size elements. If
// size ≤ 0 the whole slice is returned as a single chunk.
func Chunk(input []string, size int) [][]string {
	if size <= 0 || len(input) == 0 {
		return [][]string{input}
	}
	var chunks [][]string
	for size < len(input) {
		input, chunks = input[size:], append(chunks, input[:size:size])
	}
	return append(chunks, input)
}

// ReverseStrings returns a new slice with the elements of input in reverse order.
func ReverseStrings(input []string) []string {
	n := len(input)
	out := make([]string, n)
	for i, s := range input {
		out[n-1-i] = s
	}
	return out
}

// IntersectStrings returns the elements that are present in both a and b.
// The result preserves the order of a and contains no duplicates.
func IntersectStrings(a, b []string) []string {
	set := make(map[string]struct{}, len(b))
	for _, s := range b {
		set[s] = struct{}{}
	}
	seen := make(map[string]struct{})
	var out []string
	for _, s := range a {
		if _, inB := set[s]; inB {
			if _, already := seen[s]; !already {
				out = append(out, s)
				seen[s] = struct{}{}
			}
		}
	}
	return out
}

// DiffStrings returns the elements of a that are not present in b, preserving
// order and removing duplicates.
func DiffStrings(a, b []string) []string {
	set := make(map[string]struct{}, len(b))
	for _, s := range b {
		set[s] = struct{}{}
	}
	seen := make(map[string]struct{})
	var out []string
	for _, s := range a {
		if _, inB := set[s]; !inB {
			if _, already := seen[s]; !already {
				out = append(out, s)
				seen[s] = struct{}{}
			}
		}
	}
	return out
}
