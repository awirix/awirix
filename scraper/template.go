package scraper

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/passport"
	"strings"
	"text/template"
)

var TemplateVideo = lo.Must(template.New("video").Parse(`local {{ .Const.Module }} = {}

--- @alias {{ .Noun.Singular }} { ['str']: string, [any]: any }
--- @alias episode { ['str']: string, [any]: any }

--- Searches for the {{ .Noun.Plural }}
--- @param query string The query to search for.
--- @return {{ .Noun.Singular }}[] shows.
function {{ .Const.Module }}.{{ .Const.Search }}(query)
	return {}
end

--- Returns a list of episodes for {{ .Noun.Singular }}.
--- @param {{ .Noun.Singular }} {{ .Noun.Singular }} The {{ .Noun.Singular }} to get episodes for.
--- @return episode[] episodes.
function {{ .Const.Module }}.{{ .Const.Episodes }}({{ .Noun.Singular }})
	return {}
end

--- Prepares an episode for watching/downloading.
--- @param episode episode The episode to prepare.
--- @return episode episode The prepared episode. 
function {{ .Const.Module }}.{{ .Const.Prepare }}(episode)
	return episode
end

--- Watches an episode.
--- @param episode episode The episode to watch.
function {{ .Const.Module }}.{{ .Const.Watch }}(episode)
	require('{{ .Const.App }}').api.watch(episode.url)
end

--- Downloads an episode.
--- @param episode episode The episode to download.
function {{ .Const.Module }}.{{ .Const.Download }}(episode)
	require('{{ .Const.App }}').api.download(episode.show, episode.url)
end

return {{ .Const.Module }}`))

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
		tmpl = TemplateVideo
		noun = &Noun{"anime", "animes"}
	case passport.DomainMovies:
		tmpl = TemplateVideo
		noun = &Noun{"movie", "movies"}
	case passport.DomainShows:
		tmpl = TemplateVideo
		noun = &Noun{"show", "shows"}
	default:
		return "", fmt.Errorf("unknown domain: %s", domain)
	}

	s := struct {
		Const struct {
			Module,
			Search,
			Episodes,
			Prepare,
			Watch,
			Download,
			App string
		}
		Noun *Noun
	}{}

	s.Const.Module = constant.ScraperModuleName
	s.Const.Search = constant.FunctionSearch
	s.Const.Episodes = constant.FunctionEpisodes
	s.Const.Prepare = constant.FunctionPrepare
	s.Const.Watch = constant.FunctionWatch
	s.Const.Download = constant.FunctionDownload
	s.Const.App = constant.App
	s.Noun = noun

	var b strings.Builder
	if err := tmpl.Execute(&b, s); err != nil {
		return "", err
	}

	return b.String(), nil
}
