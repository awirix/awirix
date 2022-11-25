package scraper

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/passport"
	"strings"
	"text/template"
)

var templateVideo = lo.Must(template.New("video").Parse(`local {{ .Module }} = { {{ .Field.HasSearch }} = true }

--- @alias {{ .Noun.Singular }} { ['str']: string, [any]: any }
--- @alias episode { ['str']: string, [any]: any }

--- Searches for the {{ .Noun.Plural }}
--- @param query string The query to search for.
--- @return {{ .Noun.Singular }}[] # The {{ .Noun.Plural }} that match the query.
function {{ .Module }}.{{ .Fn.Search }}(query)
   return {}
end

--- Returns a list of episodes for {{ .Noun.Singular }}.
--- @param {{ .Noun.Singular }} {{ .Noun.Singular }}? The {{ .Noun.Singular }} to get episodes for.
--- @return episode[] # The list of episodes.
function {{ .Module }}.{{ .Fn.Episodes }}({{ .Noun.Singular }})
   return {}
end

--- Prepares an episode for watching/downloading.
--- @param episode episode The episode to prepare.
--- @return episode # The prepared episode.
function {{ .Module }}.{{ .Fn.Prepare }}(episode)
   return episode
end

--- Watches an episode.
--- @param episode episode The episode to watch.
function {{ .Module }}.{{ .Fn.Watch }}(episode)
   require('{{ .App }}').api.watch(episode.url)
end

--- Downloads an episode.
--- @param episode episode The episode to download.
function {{ .Module }}.{{ .Fn.Download }}(episode)
   require('{{ .App }}').api.download(episode.show, episode.url)
end

return {{ .Module }}

-- vim:ts=3 ss=3 sw=3 expandtab`))

func GenerateTemplate(domain passport.Domain) (string, error) {
	type Noun struct {
		Singular, Plural string
	}

	var (
		tmpl *template.Template
		noun *Noun
	)

	switch domain {
	case passport.DomainAnime:
		tmpl = templateVideo
		noun = &Noun{"anime", "animes"}
	case passport.DomainMovies:
		tmpl = templateVideo
		noun = &Noun{"movie", "movies"}
	case passport.DomainShows:
		tmpl = templateVideo
		noun = &Noun{"show", "shows"}
	default:
		return "", fmt.Errorf("unknown domain: %s", domain)
	}

	s := struct {
		Fn struct {
			Search,
			Episodes,
			Prepare,
			Watch,
			Download string
		}
		Field struct {
			HasSearch string
		}
		Module, App string
		Noun        *Noun
	}{}

	s.Module = constant.ModuleScraper
	s.App = constant.App
	s.Noun = noun
	s.Fn.Search = constant.FunctionSearch
	s.Fn.Episodes = constant.FunctionEpisodes
	s.Fn.Prepare = constant.FunctionPrepare
	s.Fn.Watch = constant.FunctionWatch
	s.Fn.Download = constant.FunctionDownload
	s.Field.HasSearch = constant.FieldHasSearch

	var b strings.Builder
	if err := tmpl.Execute(&b, s); err != nil {
		return "", err
	}

	return b.String(), nil
}
