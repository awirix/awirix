package passport

import (
	"bytes"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/github"
	"github.com/vivi-app/vivi/language"
	"github.com/vivi-app/vivi/semver"
	"github.com/vivi-app/vivi/style"
	"github.com/vivi-app/vivi/util"
	"io"
	"path"
	"strings"
	"text/template"
)

type Passport struct {
	Name         string                    `toml:"name"`
	ID           string                    `toml:"id"`
	About        string                    `toml:"about,omitempty"`
	Version      *semver.Version           `toml:"version"`
	Language     *language.Language        `toml:"lang,omitempty"`
	NSFW         bool                      `toml:"nsfw"`
	Tags         []string                  `toml:"tags,omitempty"`
	Repository   *github.Repository        `toml:"repo,omitempty"`
	Requirements *Requirements             `toml:"reqs,omitempty"`
	Config       map[string]*ConfigSection `toml:"cfg,omitempty"`
}

var passportTemplate = lo.Must(template.New("passport").Funcs(template.FuncMap{
	"join": func(s []string) string {
		return strings.Join(s, ", ")
	},
	"purple": style.Fg(color.Purple),
	"faint":  style.Faint,
	"yellow": style.Fg(color.Yellow),
	"red":    style.Fg(color.Red),
	"green":  style.Fg(color.Green),
	"bold":   style.Bold,
	"url":    style.New().Foreground(color.Cyan).Underline(true).Render,
}).Parse(`{{ purple .Name }} {{ bold .Version.String }} {{ if .NSFW }}{{ red "NSFW" }}{{ end }} 
{{ if not .About }}No description{{ else }}{{ faint .About }}{{ end }}

{{ yellow .Language.NativeName }}
{{ if not .Requirements }}No requirements{{ else }}{{ .Requirements.Info }}{{ end }}
{{ if not .Repository }}No repository{{ else }}{{ url .Repository.URL }}{{ end }}`))

func (p *Passport) Info() string {
	var b strings.Builder

	lo.Must0(passportTemplate.Execute(&b, p))

	return b.String()
}

func (p *Passport) Validate() error {
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
	err := toml.NewDecoder(reader).Decode(&passport)
	if err != nil {
		return nil, err
	}

	return &passport, passport.Validate()
}

func (p *Passport) LatestVersion() (*semver.Version, error) {
	file, err := p.Repository.GetFile(p.remotePassportPath())
	if err != nil {
		return nil, err
	}

	data, err := file.Contents()
	if err != nil {
		return nil, err
	}

	passport, err := New(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return passport.Version, nil
}

func (p *Passport) remotePath() string {
	return path.Join(constant.ExtensionsDir, util.SanitizeFilename(p.ID))
}

func (p *Passport) remotePassportPath() string {
	return path.Join(p.remotePath(), constant.FilenamePassport)
}
