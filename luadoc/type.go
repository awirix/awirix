package luadoc

import (
	"fmt"
	"github.com/samber/lo"
	"strings"
)

const (
	String  = "string"
	Number  = "number"
	Boolean = "boolean"
	Table   = "table"
	Any     = "any"
)

func List(of string) string {
	return of + "[]"
}

func Map(keys string, values string) string {
	return Table + "<" + keys + ", " + values + ">"
}

func Enum(members ...string) string {
	quoted := lo.Map(members, func(s string, _ int) string {
		return fmt.Sprintf(`'%s'`, s)
	})

	return strings.Join(quoted, " | ")
}
