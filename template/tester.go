package template

import (
	"bytes"
	_ "embed"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/tester"
	"text/template"
)

//go:embed tester.lua.tmpl
var templateLuaTester string

func LuaTester() []byte {
	tmpl := lo.Must(template.New(tester.Module).Parse(templateLuaTester))
	m := newMeta(tester.Module)

	var b bytes.Buffer

	lo.Must0(tmpl.Execute(&b, m))

	return b.Bytes()
}

//go:embed tester.fnl.tmpl
var templateFennelTester string

func FennelTester() []byte {
	tmpl := lo.Must(template.New(tester.Module).Parse(templateFennelTester))
	m := newMeta(tester.Module)

	var b bytes.Buffer

	lo.Must0(tmpl.Execute(&b, m))

	return b.Bytes()
}
