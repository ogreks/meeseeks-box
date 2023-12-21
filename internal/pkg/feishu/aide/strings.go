package aide

import "strings"

// CutPrefix prefix matches
func CutPrefix(s, prefix string) (string, bool) {
	if strings.HasPrefix(s, prefix) {
		return strings.TrimPrefix(s, prefix), true
	}

	return s, false
}

// EitherCutPrefix If any prefix matches, the remainder is returned
func EitherCutPrefix(s string, prefixs ...string) (string, bool) {
	for _, p := range prefixs {
		if strings.HasPrefix(s, p) {
			return strings.TrimPrefix(s, p), true
		}
	}

	return s, false
}

// TrimEqual trim space and equal
func TrimEqual(s, prefix string) (string, bool) {
	if strings.TrimSpace(s) == prefix {
		return "", true
	}

	return s, false
}

// EitherTrimEqual
func EitherTrimEqual(s string, prefixs ...string) (string, bool) {
	for _, p := range prefixs {
		if strings.TrimSpace(s) == p {
			return "", true
		}
	}

	return s, false
}
