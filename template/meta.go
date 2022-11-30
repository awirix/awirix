package template

import "github.com/vivi-app/vivi/constant"

type funcs struct {
	Search,
	Explore,
	Prepare,
	Play,
	Download,
	Test string
}

type meta struct {
	Module string
	App    string
	Fn     *funcs
}

func newMeta(module string) *meta {
	m := &meta{}

	m.Module = module
	m.App = constant.App
	m.Fn = &funcs{
		Search:   constant.FunctionSearch,
		Explore:  constant.FunctionExplore,
		Prepare:  constant.FunctionPrepare,
		Play:     constant.FunctionPlay,
		Download: constant.FunctionDownload,
		Test:     constant.FunctionTest,
	}

	return m
}
