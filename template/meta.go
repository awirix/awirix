package template

import (
	"github.com/vivi-app/vivi/app"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/tester"
)

type funcs struct {
	Search,
	Explore,
	Prepare,
	Play,
	Download,
	Test string
}

type fields struct {
	Display string
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
	}
	m.Fn = &funcs{
		Search:   scraper.FunctionSearch,
		Explore:  scraper.FunctionExplore,
		Prepare:  scraper.FunctionPrepare,
		Play:     scraper.FunctionPlay,
		Download: scraper.FunctionDownload,
		Test:     tester.FunctionTest,
	}

	return m
}
