package template

import (
	"bytes"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/constant"
	"text/template"
)

func NewTest() []byte {
	tmpl := lo.Must(template.New("test").
		Parse(`-- vim:ts=3 ss=3 sw=3 expandtab

local M = {}

function M.{{ .Fn.Test }}()
	assert(2 + 2 == 4, 'Math is broken')
end

return M`))

	m := newMeta(constant.ModuleTest, nil)

	var b bytes.Buffer

	lo.Must0(tmpl.Execute(&b, m))

	return b.Bytes()
}
