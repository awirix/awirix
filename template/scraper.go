package template

import (
	"bytes"
	_ "embed"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/scraper"
	"text/template"
)

//go:embed scraper.lua.tmpl
var templateLuaScraper string

func LuaScraper() []byte {
	tmpl := lo.Must(template.New(scraper.Module).Parse(templateLuaScraper))
	info := newMeta(scraper.Module)

	var b bytes.Buffer
	lo.Must0(tmpl.Execute(&b, info))

	return b.Bytes()
}

//go:embed scraper.fnl.tmpl
var templateFennelScraper string

func FennelScraper() []byte {
	tmpl := lo.Must(template.New(scraper.Module).Parse(templateFennelScraper))
	info := newMeta(scraper.Module)

	var b bytes.Buffer
	lo.Must0(tmpl.Execute(&b, info))

	return b.Bytes()
}
