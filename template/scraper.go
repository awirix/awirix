package template

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/passport"
	"text/template"
)

//go:embed scraper.lua.tmpl
var templateVideo string

func NewScraper(domain passport.Domain) []byte {
	var (
		tmpl *template.Template
		n    *noun
	)

	switch domain {
	case passport.DomainAnime:
		tmpl = lo.Must(template.New(constant.ModuleScraper).Parse(templateVideo))
		n = &noun{"anime", "animes"}
	case passport.DomainMovies:
		tmpl = lo.Must(template.New(constant.ModuleScraper).Parse(templateVideo))
		n = &noun{"movie", "movies"}
	case passport.DomainShows:
		tmpl = lo.Must(template.New(constant.ModuleScraper).Parse(templateVideo))
		n = &noun{"show", "shows"}
	default:
		panic(fmt.Sprintf("unknown domain: %s", domain))
	}

	info := newMeta(constant.ModuleScraper, n)

	var b bytes.Buffer
	lo.Must0(tmpl.Execute(&b, info))

	return b.Bytes()
}
