package filename

import (
	"github.com/samber/lo"
	"regexp"
)

// replacers is a list of regexp.Regexp pairs that will be used to sanitize filenames.
var sanitizeReplacers = []lo.Tuple2[*regexp.Regexp, string]{
	{regexp.MustCompile(`[\\/<>:;"'|?!*{}#%&^+,~\s]`), "_"},
	{regexp.MustCompile(`__+`), "_"},
	{regexp.MustCompile(`^[_\-.]+|[_\-.]+$`), ""},
}

// Sanitize will remove all invalid characters from a path.
func Sanitize(filename string) string {
	for _, re := range sanitizeReplacers {
		filename = re.A.ReplaceAllString(filename, re.B)
	}

	return filename
}
