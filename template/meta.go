package template

import (
	"github.com/vivi-app/vivi/app"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/tester"
)

type funcs struct {
	Search,
	Prepare,
	Stream,
	Download,
	Test string
}

type fields struct {
	Display string
	About   string
	Layers  string
}

type meta struct {
	Module string
	Fields *fields
	App    string
	Fn     *funcs
}

func newMeta(module string) *meta {
	m := &meta{}

	m.Module = module
	m.App = app.Name
	m.Fields = &fields{
		Display: scraper.FieldDisplay,
		About:   scraper.FieldAbout,
	}
	m.Fn = &funcs{
		Search:   scraper.FunctionSearch,
		Prepare:  scraper.FunctionPrepare,
		Stream:   scraper.FunctionStream,
		Download: scraper.FunctionDownload,
		Test:     tester.FunctionTest,
	}

	return m
}
