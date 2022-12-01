package template

import (
	"bytes"
	_ "embed"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/scraper"
	"text/template"
)

//go:embed scraper.lua.tmpl
var templateScraper string

func Scraper() []byte {
	tmpl := lo.Must(template.New(scraper.Module).Parse(templateScraper))

	info := newMeta(scraper.Module)

	var b bytes.Buffer
	lo.Must0(tmpl.Execute(&b, info))

	return b.Bytes()
}
