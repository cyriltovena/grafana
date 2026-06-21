package util

// MergeStringMaps merges any number of string→string maps left-to-right.
// Keys in later maps overwrite keys in earlier maps. A new map is always
// returned; the inputs are not modified.
func MergeStringMaps(maps ...map[string]string) map[string]string {
	out := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// CopyStringMap returns a shallow copy of m.
func CopyStringMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// StringMapKeys returns the keys of m in an unspecified order.
func StringMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// StringMapValues returns the values of m in an unspecified order.
func StringMapValues(m map[string]string) []string {
	vals := make([]string, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

// FilterStringMap returns a new map containing only the key-value pairs for
// which keep returns true.
func FilterStringMap(m map[string]string, keep func(k, v string) bool) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		if keep(k, v) {
			out[k] = v
		}
	}
	return out
}

// InvertStringMap returns a new map with keys and values swapped. If m
// contains duplicate values only one of the corresponding keys will survive
// (which one is unspecified).
func InvertStringMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[v] = k
	}
	return out
}

// MapFromSlice builds a string→string map from a flat slice of alternating
// key-value pairs. If the slice has an odd length the last element is silently
// ignored.
func MapFromSlice(kvs []string) map[string]string {
	out := make(map[string]string, len(kvs)/2)
	for i := 0; i+1 < len(kvs); i += 2 {
		out[kvs[i]] = kvs[i+1]
	}
	return out
}
