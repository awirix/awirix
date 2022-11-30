package template

import (
	"bytes"
	_ "embed"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/constant"
	"text/template"
)

//go:embed test.lua.tmpl
var templateTest string

func NewTest() []byte {
	tmpl := lo.Must(template.New("test").Parse(templateTest))

	m := newMeta(constant.ModuleTest)

	var b bytes.Buffer

	lo.Must0(tmpl.Execute(&b, m))

	return b.Bytes()
}
