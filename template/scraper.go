package template

import (
	"bytes"
	"fmt"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/passport"
	"text/template"
)

var templateVideo = lo.Must(template.New("video").Parse(`-- vim:ts=3 ss=3 sw=3 expandtab

--- @alias stringer { ['string']: string, [any]: any }

local M = {}

--- Searches for the {{ .Noun.Plural }}
--- @param query string The query to search for.
--- @return stringer[] # The {{ .Noun.Plural }} that match the query.
function M.{{ .Fn.Search }}(query)
   return {}
end

--- Returns a list of episodes for {{ .Noun.Singular }}.
--- @param {{ .Noun.Singular }} stringer? The {{ .Noun.Singular }} to get episodes for.
--- @return stringer[] # The list of episodes.
function M.{{ .Fn.Episodes }}({{ .Noun.Singular }})
   return {}
end

--- Prepares an episode for watching/downloading.
--- @param episode stringer The episode to prepare.
--- @return stringer # The prepared episode.
function M.{{ .Fn.Prepare }}(episode)
   return episode
end

--- Watches an episode.
--- @param episode stringer The episode to watch.
function M.{{ .Fn.Watch }}(episode)
   require('{{ .App }}').api.watch(episode.url)
end

--- Downloads an episode.
--- @param episode stringer The episode to download.
function M.{{ .Fn.Download }}(episode)
   require('{{ .App }}').api.download(episode.string, episode.url)
end

return M
`))

func NewScraper(domain passport.Domain) []byte {
	var (
		tmpl *template.Template
		n    *noun
	)

	switch domain {
	case passport.DomainAnime:
		tmpl = templateVideo
		n = &noun{"anime", "animes"}
	case passport.DomainMovies:
		tmpl = templateVideo
		n = &noun{"movie", "movies"}
	case passport.DomainShows:
		tmpl = templateVideo
		n = &noun{"show", "shows"}
	default:
		panic(fmt.Sprintf("unknown domain: %s", domain))
	}

	info := newMeta(constant.ModuleScraper, n)

	var b bytes.Buffer
	lo.Must0(tmpl.Execute(&b, info))

	return b.Bytes()
}
