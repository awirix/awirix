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
	"github.com/enescakir/emoji"
	"github.com/samber/lo"
	"io"
	"strings"
	"text/template"
)

type Passport struct {
	Icon         emoji.Emoji        `json:"-"`
	Name         string             `json:"name"`
	ID           string             `json:"id"`
	About        string             `json:"about"`
	Version      *version.Version   `json:"version"`
	Language     *language.Language `json:"language"`
	NSFW         bool               `json:"nsfw"`
	Tags         []string           `json:"tags,omitempty"`
	Repository   *github.Repository `json:"repository,omitempty"`
	Requirements *Requirements      `json:"requirements,omitempty"`
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

func (p *Passport) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("passport: id is empty")
	}

	if p.Name == "" {
		return fmt.Errorf("passport: name is empty")
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
	var passport = &Passport{}

	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(passport)
	if err != nil {
		return nil, err
	}

	return passport, passport.Validate()
}
