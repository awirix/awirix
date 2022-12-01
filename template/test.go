package template

import (
	"bytes"
	_ "embed"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/tester"
	"text/template"
)

//go:embed test.lua.tmpl
var templateTester string

func Tester() []byte {
	tmpl := lo.Must(template.New(tester.Module).Parse(templateTester))

	m := newMeta(tester.Module)

	var b bytes.Buffer

	lo.Must0(tmpl.Execute(&b, m))

	return b.Bytes()
}
