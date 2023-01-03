package luadoc

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
