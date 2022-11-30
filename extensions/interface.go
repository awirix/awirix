package extensions

import "github.com/vivi-app/vivi/extensions/passport"

type ExtensionContainer interface {
	String() string
	Path() string
	Passport() *passport.Passport
}
