package template

import "github.com/vivi-app/vivi/constant"

type noun struct {
	Singular, Plural string
}

type funcs struct {
	Search,
	Episodes,
	Prepare,
	Watch,
	Download,
	Test string
}

type meta struct {
	Noun   *noun
	Module string
	App    string
	Fn     *funcs
}

func newMeta(module string, n *noun) *meta {
	m := &meta{}

	m.Noun = n
	m.Module = module
	m.App = constant.App
	m.Fn = &funcs{
		Search:   constant.FunctionSearch,
		Episodes: constant.FunctionEpisodes,
		Prepare:  constant.FunctionPrepare,
		Watch:    constant.FunctionWatch,
		Download: constant.FunctionDownload,
		Test:     constant.FunctionTest,
	}

	return m
}
