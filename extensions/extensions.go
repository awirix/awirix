package extensions

import "github.com/awirix/awirix/extensions/passport"

type ExtensionWrapper interface {
	String() string
	Path() string
	Cache() string
	Downloads() string
	Passport() *passport.Passport
}
