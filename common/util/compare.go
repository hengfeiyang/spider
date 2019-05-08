package util

// Eq ...
func Eq(v1, v2 interface{}) bool {
	return v1 == v2
}

// Neq ...
func Neq(v1, v2 interface{}) bool {
	return v1 != v2
}

// Gt ...
func Gt(v1, v2 int64) bool {
	return v1 > v2
}

// Lt ...
func Lt(v1, v2 int64) bool {
	return v1 < v2
}

// Empty ...
func Empty(v interface{}) bool {
	if v == nil {
		return true
	}

	switch v.(type) {
	case string:
		if v == "" {
			return true
		}
	case int:
		if v == 0 {
			return true
		}
	case int32:
		if v == int32(0) {
			return true
		}
	case int64:
		if v == int64(0) {
			return true
		}
	case float32:
		if v == float32(0.0) {
			return true
		}
	case float64:
		if v == float64(0.0) {
			return true
		}
	case bool:
		if v == false {
			return true
		}
	}

	return false
}
