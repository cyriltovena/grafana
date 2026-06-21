package util

import (
	"reflect"
	"sort"
	"testing"
)

func TestContainsString(t *testing.T) {
	s := []string{"a", "b", "c"}
	if !ContainsString(s, "b") {
		t.Error("ContainsString: expected true for 'b'")
	}
	if ContainsString(s, "z") {
		t.Error("ContainsString: expected false for 'z'")
	}
	if ContainsString(nil, "a") {
		t.Error("ContainsString: expected false for nil slice")
	}
}

func TestContainsInt(t *testing.T) {
	s := []int{1, 2, 3}
	if !ContainsInt(s, 2) {
		t.Error("ContainsInt: expected true for 2")
	}
	if ContainsInt(s, 99) {
		t.Error("ContainsInt: expected false for 99")
	}
}

func TestUniqueStrings(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b"}
	got := UniqueStrings(input)
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("UniqueStrings = %v; want %v", got, want)
	}
}

func TestFilterStrings(t *testing.T) {
	input := []string{"foo", "bar", "baz", "qux"}
	got := FilterStrings(input, func(s string) bool { return len(s) == 3 })
	want := []string{"foo", "bar", "baz", "qux"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FilterStrings = %v; want %v", got, want)
	}

	got2 := FilterStrings(input, func(s string) bool { return s[0] == 'b' })
	want2 := []string{"bar", "baz"}
	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("FilterStrings = %v; want %v", got2, want2)
	}
}

func TestMapStrings(t *testing.T) {
	input := []string{"hello", "world"}
	got := MapStrings(input, func(s string) string { return s + "!" })
	want := []string{"hello!", "world!"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MapStrings = %v; want %v", got, want)
	}
}

func TestChunk(t *testing.T) {
	input := []string{"a", "b", "c", "d", "e"}
	got := Chunk(input, 2)
	if len(got) != 3 {
		t.Fatalf("Chunk(5, 2) returned %d chunks; want 3", len(got))
	}
	if !reflect.DeepEqual(got[0], []string{"a", "b"}) {
		t.Errorf("chunk[0] = %v; want [a b]", got[0])
	}
	if !reflect.DeepEqual(got[2], []string{"e"}) {
		t.Errorf("chunk[2] = %v; want [e]", got[2])
	}

	// size 0 should return one chunk.
	if c := Chunk(input, 0); len(c) != 1 {
		t.Errorf("Chunk(size=0) returned %d chunks; want 1", len(c))
	}
}

func TestReverseStrings(t *testing.T) {
	got := ReverseStrings([]string{"a", "b", "c"})
	want := []string{"c", "b", "a"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReverseStrings = %v; want %v", got, want)
	}
}

func TestIntersectStrings(t *testing.T) {
	a := []string{"x", "y", "z", "y"}
	b := []string{"y", "w", "z"}
	got := IntersectStrings(a, b)
	sort.Strings(got)
	want := []string{"y", "z"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("IntersectStrings = %v; want %v", got, want)
	}
}

func TestDiffStrings(t *testing.T) {
	a := []string{"a", "b", "c", "d"}
	b := []string{"b", "d"}
	got := DiffStrings(a, b)
	want := []string{"a", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DiffStrings = %v; want %v", got, want)
	}
}
