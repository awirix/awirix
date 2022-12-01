package extensions

import "github.com/vivi-app/vivi/extensions/passport"

const Dirname = "extensions"

type ExtensionContainer interface {
	String() string
	Path() string
	Passport() *passport.Passport
}
