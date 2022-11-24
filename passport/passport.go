package passport

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/icon"
	"github.com/vivi-app/vivi/language"
	"github.com/vivi-app/vivi/semver"
	"github.com/vivi-app/vivi/style"
	"strings"
	"text/template"
)

type Passport struct {
	Name         string                    `toml:"name"`
	ID           string                    `toml:"id"`
	About        string                    `toml:"about,omitempty"`
	Version      semver.Version            `toml:"version"`
	Domain       Domain                    `toml:"domain"`
	Language     language.Language         `toml:"lang,omitempty"`
	NSFW         bool                      `toml:"nsfw"`
	Tags         []string                  `toml:"tags,omitempty"`
	Dependencies Requirements              `toml:"reqs,omitempty"`
	Config       map[string]*ConfigSection `toml:"cfg,omitempty"`
}

var passportTemplate = lo.Must(template.New("passport").Funcs(template.FuncMap{
	"join": func(s []string) string {
		return strings.Join(s, ", ")
	},
	"purple":    style.Fg(color.Purple),
	"faint":     style.Faint,
	"red":       style.Fg(color.Red),
	"green":     style.Fg(color.Green),
	"bold":      style.Bold,
	"iconCheck": func() string { return icon.Check },
	"iconCross": func() string { return icon.Cross },
}).Parse(`{{ purple .Name }} {{ bold .Version.String }} {{ if .NSFW }}{{ red "NSFW" }}{{ end }}
{{ faint .About }}

{{ .Requirements.Info }}`))

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

	return nil
}

func (p *Passport) CheckDependencies() bool {
	return p.Dependencies.Matches()
}
