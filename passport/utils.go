package passport

import "strings"

func ToID(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "-"))
}
