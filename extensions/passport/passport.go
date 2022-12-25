package passport

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/github"
	"github.com/vivi-app/vivi/icon"
	"github.com/vivi-app/vivi/language"
	"github.com/vivi-app/vivi/style"
	"github.com/vivi-app/vivi/version"
	"io"
	"strings"
	"text/template"
)

type Passport struct {
	MetaSchema   string                    `json:"$schema,omitempty"`
	Name         string                    `json:"name" jsonschema:"description=title of the extension"`
	ID           string                    `json:"id" jsonschema:"description=ID of the extension"`
	About        string                    `json:"about" jsonschema:"description=Brief description"`
	VersionRaw   string                    `json:"version" jsonschema:"required,description=Semantic version of the extension,example=1.0.0"`
	LanguageRaw  string                    `json:"language" jsonschema:"required,description=Primary language of the content"`
	NSFW         bool                      `json:"nsfw" jsonschema:"description=Whether the extension contains NSFW content"`
	Tags         []string                  `json:"tags,omitempty" jsonschema:"description=Tags"`
	Repository   *github.Repository        `json:"repository,omitempty" jsonschema:"description=Repository which contains the extension"`
	Requirements *Requirements             `json:"requirements,omitempty" jsonschema:"description=Requirements that must be met to use the extension"`
	Config       map[string]*ConfigSection `json:"config,omitempty" jsonschema:"description=Configuration fields"`
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
{{ if .Requirements }}{{ .Requirements.Info }}{{ end }}
{{ if .Repository }}{{ url .Repository.URL }}{{ end }}`))

func (p *Passport) Info() string {
	var b strings.Builder

	lo.Must0(passportTemplate.Execute(&b, p))

	return strings.TrimSpace(b.String())
}

func (p *Passport) Version() *version.Version {
	return version.MustParse(p.VersionRaw)
}

func (p *Passport) Language() *language.Language {
	return language.Languages[p.LanguageRaw]
}

func (p *Passport) Validate() error {
	if _, err := version.NewVersion(p.VersionRaw); err != nil {
		return fmt.Errorf("invalid version: %s", err)
	}

	if _, ok := language.Languages[p.LanguageRaw]; !ok {
		return fmt.Errorf("invalid ISO 639-1 language code: %s", p.LanguageRaw)
	}

	if p.Repository != nil {
		for _, t := range []*lo.Tuple2[string, string]{
			{"name", p.Repository.Name},
			{"owner", p.Repository.Owner},
		} {
			if t.B == "" {
				return fmt.Errorf("missing required field in repo: %s", t.A)
			}
		}
	}

	return nil
}

func (p *Passport) CheckRequirements() bool {
	return p.Requirements.Matches()
}

func New(reader io.Reader) (*Passport, error) {
	var passport Passport
	err := json.NewDecoder(reader).Decode(&passport)
	if err != nil {
		return nil, err
	}

	return &passport, passport.Validate()
}
