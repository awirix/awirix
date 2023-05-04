package passport

import (
	"encoding/json"
	"fmt"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/github"
	"github.com/awirix/awirix/icon"
	"github.com/awirix/awirix/language"
	"github.com/awirix/awirix/style"
	"github.com/awirix/awirix/version"
	"github.com/mvdan/xurls"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"io"
	"net/url"
	"strings"
	"text/template"
)

type Passport struct {
	Name       string                       `json:"name"`
	ID         string                       `json:"id"`
	About      string                       `json:"about"`
	Version    version.Version              `json:"version"`
	Awirix     version.Version              `json:"awirix"`
	Language   language.Language            `json:"language"`
	Authors    []string                     `json:"authors"`
	NSFW       bool                         `json:"nsfw"`
	Tags       []string                     `json:"tags"`
	Repository mo.Option[github.Repository] `json:"repository"`
	Programs   []string                     `json:"programs"`
}

func errPassport(err error) error {
	return errors.Wrap(err, "passport")
}

var passportTemplate = lo.Must(template.New("passport").Funcs(template.FuncMap{
	"tags": func(tags []string) string {
		if len(tags) > 5 {
			tags = tags[:5]
		}

		return strings.Join(tags, style.Fg(color.Blue)(" "+icon.CDot+" "))
	},
	"nsfw":   style.NewColored(color.New("#fff"), color.New("#ff0000")).Bold(true).Padding(0, 1).Render,
	"purple": style.Fg(color.Purple),
	"faint":  style.Faint,
	"yellow": style.Fg(color.Yellow),
	"red":    style.Fg(color.Red),
	"green":  style.Fg(color.Green),
	"bold":   style.Bold,
	"cyan":   style.Fg(color.Cyan),
	"url":    style.New().Foreground(color.Cyan).Underline(true).Render,
}).Parse(`{{ bold (purple .Name) }} {{ bold .Version.String }} {{ if .NSFW }}{{ nsfw "NSFW" }}{{ end }} 
{{ if not .About }}{{ faint "No description" }}{{ else }}{{ faint .About }}{{ end }}

{{ yellow (.Language).Name }} {{ if not (eq (.Language).Name (.Language).NativeName) }}{{ faint (.Language).NativeName }}{{ end }}
{{ if .Tags }}{{ tags .Tags }}{{ end }}
{{ if .Repository }}{{ url .Repository.URL }}{{ end }}`))

func (p *Passport) Info() string {
	var b strings.Builder

	lo.Must0(passportTemplate.Execute(&b, p))

	return strings.TrimSpace(b.String())
}

func (p *Passport) InfoMarkdown() string {
	tmpl := template.Must(template.New("passport").Funcs(map[string]any{
		"code": func(code string) string {
			return fmt.Sprintf("`%s`", code)
		},
		"format_urls": func(s string) string {
			return xurls.Relaxed.ReplaceAllStringFunc(s, func(rawURL string) string {
				if !strings.HasPrefix(rawURL, "https://") || !strings.HasPrefix(rawURL, "http://") {
					rawURL = "https://" + rawURL
				}

				u, err := url.Parse(rawURL)
				if err != nil {
					return fmt.Sprintf("[%[1]s](%[1]s)", rawURL)
				}

				return fmt.Sprintf("[%s](%s)", u.Hostname(), u.String())
			})
		},
	}).Parse(`
{{- "##" }} {{ .Name }} {{ .Version.String }}

ID {{ code .ID }}

{{ if .NSFW }}This extension is marked as NSFW{{ end }}

Primary language is **{{ .Language.Name }}**

{{ if .Programs }}Programs required{{ end }}
{{ range $program := .Programs }}
- {{ $program }}
{{ end }}

{{ if .Repository }}[GitHub Repository]({{ .Repository.URL }}){{ end }}

### About

{{ if .Tags }}**Tags**{{ end }}
{{ range $tag := .Tags }}
- {{ $tag }}
{{ end }}

> {{ format_urls .About }}
`))

	var b strings.Builder
	lo.Must0(tmpl.Execute(&b, p))

	return strings.TrimSpace(b.String())
}

func (p *Passport) Validate() error {
	if p.ID == "" {
		return errPassport(fmt.Errorf("id is empty"))
	}

	if p.Name == "" {
		return errPassport(fmt.Errorf("name is empty"))
	}

	if repo, ok := p.Repository.Get(); ok {
		for _, t := range []*lo.Tuple2[string, string]{
			{"name", repo.Name},
			{"owner", repo.Owner},
		} {
			if t.B == "" {
				return errPassport(fmt.Errorf("missing required field in repo: %s", t.A))
			}
		}
	}

	return nil
}

func New(reader io.Reader) (*Passport, error) {
	var passport = &Passport{}

	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(passport)
	if err != nil {
		return nil, err
	}

	return passport, passport.Validate()
}
