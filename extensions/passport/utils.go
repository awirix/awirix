package passport

import "strings"

func ToID(s string) string {
	// allow only lowercase alphanumeric characters without spaces
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r
		}

		if r >= 'A' && r <= 'Z' {
			// to lowercase
			return r + 32
		}

		if r >= '0' && r <= '9' {
			return r
		}

		if r == ' ' {
			return '-'
		}

		return -1
	}, s)
}
