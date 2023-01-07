package extensions

import "github.com/awirix/awirix/extensions/passport"

type ExtensionContainer interface {
	String() string
	Path() string
	Downloads() string
	Passport() *passport.Passport
}
