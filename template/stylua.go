package template

import _ "embed"

//go:embed stylua.toml.tmpl
var styluaTemplate string

func Stylua() []byte {
	return []byte(styluaTemplate)
}