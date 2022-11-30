package template

import (
	"bytes"
	_ "embed"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/constant"
	"text/template"
)

//go:embed scraper.lua.tmpl
var templateVideo string

func NewScraper() []byte {
	tmpl := lo.Must(template.New(constant.ModuleScraper).Parse(templateVideo))

	info := newMeta(constant.ModuleScraper)

	var b bytes.Buffer
	lo.Must0(tmpl.Execute(&b, info))

	return b.Bytes()
}
