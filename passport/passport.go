package passport

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/language"
	"github.com/vivi-app/vivi/semver"
	"github.com/vivi-app/vivi/style"
	"io"
	"path/filepath"
	"strings"
	"text/template"
)

type Passport struct {
	Name         string                    `toml:"name"`
	ID           string                    `toml:"id"`
	About        string                    `toml:"about,omitempty"`
	Version      *semver.Version           `toml:"version"`
	Domain       Domain                    `toml:"domain"`
	Language     *language.Language        `toml:"lang,omitempty"`
	NSFW         bool                      `toml:"nsfw"`
	Tags         []string                  `toml:"tags,omitempty"`
	Github       *Github                   `toml:"github,omitempty"`
	Requirements *Requirements             `toml:"reqs,omitempty"`
	Config       map[string]*ConfigSection `toml:"cfg,omitempty"`
}

var passportTemplate = lo.Must(template.New("passport").Funcs(template.FuncMap{
	"join": func(s []string) string {
		return strings.Join(s, ", ")
	},
	"purple": style.Fg(color.Purple),
	"faint":  style.Faint,
	"red":    style.Fg(color.Red),
	"green":  style.Fg(color.Green),
	"bold":   style.Bold,
	"url":    style.New().Foreground(color.Cyan).Underline(true).Render,
}).Parse(`{{ purple .Name }} {{ bold .Version.String }} {{ if .NSFW }}{{ red "NSFW" }}{{ end }}
{{ if not .About }}No description{{ else }}{{ faint .About }}{{ end }}

{{ if not .Requirements}}No requirements{{ else }}{{ .Requirements.Info }}{{ end }}
{{ if not .Github }}No repository{{ else }}{{ url .Github.URL }}{{ end }}`))

func (p *Passport) String() string {
	var b strings.Builder

	lo.Must0(passportTemplate.Execute(&b, p))

	return b.String()
}

func (p *Passport) Validate() error {
	// check if domain is valid
	if !lo.Contains(Domains, p.Domain) {
		return fmt.Errorf("invalid domain: %s", p.Domain)
	}

	if p.Github != nil {
		if err := p.Github.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Passport) CheckDependencies() bool {
	return p.Requirements.Matches()
}

func Parse(reader io.Reader) (*Passport, error) {
	var passport Passport
	err := toml.NewDecoder(reader).Decode(&passport)
	if err != nil {
		return nil, err
	}

	return &passport, passport.Validate()
}

func FromPath(path string) (*Passport, error) {
	isDir, err := filesystem.Api().IsDir(path)
	if err != nil {
		return nil, err
	}

	if isDir {
		path = filepath.Join(path, constant.Passport)
	}

	file, err := filesystem.Api().Open(path)
	if err != nil {
		return nil, err
	}

	passport, err := Parse(file)
	if err != nil {
		return nil, err
	}

	return passport, nil
}
