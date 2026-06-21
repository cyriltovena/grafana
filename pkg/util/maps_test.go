package util

import (
	"reflect"
	"sort"
	"testing"
)

func TestMergeStringMaps(t *testing.T) {
	a := map[string]string{"k1": "v1", "k2": "v2"}
	b := map[string]string{"k2": "override", "k3": "v3"}
	got := MergeStringMaps(a, b)
	want := map[string]string{"k1": "v1", "k2": "override", "k3": "v3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergeStringMaps = %v; want %v", got, want)
	}
	// Inputs should be unchanged.
	if a["k2"] != "v2" {
		t.Error("MergeStringMaps mutated input a")
	}

	// Merging zero maps returns empty map.
	if m := MergeStringMaps(); len(m) != 0 {
		t.Errorf("MergeStringMaps() = %v; want empty", m)
	}
}

func TestCopyStringMap(t *testing.T) {
	orig := map[string]string{"a": "1", "b": "2"}
	cp := CopyStringMap(orig)
	if !reflect.DeepEqual(orig, cp) {
		t.Errorf("CopyStringMap content mismatch")
	}
	cp["c"] = "3"
	if _, ok := orig["c"]; ok {
		t.Error("CopyStringMap: mutation of copy affected original")
	}
}

func TestStringMapKeys(t *testing.T) {
	m := map[string]string{"x": "1", "y": "2", "z": "3"}
	keys := StringMapKeys(m)
	sort.Strings(keys)
	if !reflect.DeepEqual(keys, []string{"x", "y", "z"}) {
		t.Errorf("StringMapKeys = %v; want [x y z]", keys)
	}
}

func TestStringMapValues(t *testing.T) {
	m := map[string]string{"a": "apple"}
	vals := StringMapValues(m)
	if len(vals) != 1 || vals[0] != "apple" {
		t.Errorf("StringMapValues = %v; want [apple]", vals)
	}
}

func TestFilterStringMap(t *testing.T) {
	m := map[string]string{"a": "apple", "b": "banana", "c": "cherry"}
	got := FilterStringMap(m, func(k, v string) bool { return k != "b" })
	if _, ok := got["b"]; ok {
		t.Error("FilterStringMap should not include 'b'")
	}
	if len(got) != 2 {
		t.Errorf("FilterStringMap len = %d; want 2", len(got))
	}
}

func TestInvertStringMap(t *testing.T) {
	m := map[string]string{"key": "value"}
	inv := InvertStringMap(m)
	if inv["value"] != "key" {
		t.Errorf("InvertStringMap = %v; want map[value:key]", inv)
	}
}

func TestMapFromSlice(t *testing.T) {
	kvs := []string{"k1", "v1", "k2", "v2"}
	got := MapFromSlice(kvs)
	want := map[string]string{"k1": "v1", "k2": "v2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MapFromSlice = %v; want %v", got, want)
	}
	// Odd-length input: last element ignored.
	odd := MapFromSlice([]string{"k1", "v1", "k2"})
	if _, ok := odd["k2"]; ok {
		t.Error("MapFromSlice with odd-length: should not have key 'k2' without a value")
	}
}
