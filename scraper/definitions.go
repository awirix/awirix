package scraper

import "github.com/vivi-app/lua"

const Module = "scraper"

const (
	FunctionPrepare  = "prepare"
	FunctionStream   = "stream"
	FunctionDownload = "download"
)

const (
	FieldSearch = "search"
	FieldLayers = "layers"

	FieldTitle       = "title"
	FieldSubtitle    = "subtitle"
	FieldPlaceholder = "placeholder"
	FieldDescription = "description"
	FieldHandler     = "handler"

	FieldNoun     = "noun"
	FieldSingular = "singular"
	FieldPlural   = "plural"
)

type Noun struct {
	singular,
	plural string
}

func (n *Noun) Singular() string {
	singular := n.singular
	if singular == "" {
		singular = "media"
	}

	return singular
}

func (n *Noun) Plural() string {
	plural := n.plural
	if plural == "" {
		plural = n.Singular() + "s"
	}

	return plural
}

type Layer struct {
	scraper *Scraper
	title   string
	handler *lua.LFunction
	noun    *Noun
}

func (l *Layer) Noun() *Noun {
	return l.noun
}

func (l *Layer) Title() string {
	if title := l.title; title != "" {
		return title
	}

	return "Select a " + l.Noun().Singular()
}

func (l *Layer) Call(media *Media) (subMedia []*Media, err error) {
	var value lua.LValue
	if media != nil {
		value = media.Value()
	} else {
		value = lua.LNil
	}

	err = l.scraper.state.CallByParam(lua.P{
		Fn:      l.handler,
		NRet:    1,
		Protect: true,
	}, value, l.scraper.progress)
	if err != nil {
		return nil, err
	}

	return l.scraper.checkMediaSlice()
}

type Search struct {
	scraper     *Scraper
	title       string
	subtitle    string
	placeholder string
	handler     *lua.LFunction
	noun        *Noun
}

func (s *Search) Title() string {
	if title := s.title; title != "" {
		return title
	}

	return "Search"
}

func (s *Search) Subtitle() string {
	if subtitle := s.subtitle; subtitle != "" {
		return subtitle
	}

	return "Select a " + s.Noun().Singular()
}

func (s *Search) Placeholder() string {
	if placeholder := s.placeholder; placeholder != "" {
		return placeholder
	}

	return "Search " + s.Noun().Plural()
}

func (s *Search) Call(query string) (subMedia []*Media, err error) {
	err = s.scraper.state.CallByParam(lua.P{
		Fn:      s.handler,
		NRet:    1,
		Protect: true,
	}, lua.LString(query), s.scraper.progress)
	if err != nil {
		return nil, err
	}

	return s.scraper.checkMediaSlice()
}

func (s *Search) Noun() *Noun {
	return s.noun
}
