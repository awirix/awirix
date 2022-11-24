package passport

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/github"
	"github.com/vivi-app/vivi/semver"
	"net/url"
	"path"
)

type Github struct {
	Repository github.Repository `toml:"repo"`
	Path       string            `toml:"path"`
}

func (g *Github) URL() string {
	res, _ := url.JoinPath(g.Repository.URL(), g.Path)
	return res
}

func (g *Github) PassportPath() string {
	return path.Join(g.Path, constant.Passport)
}

func (g *Github) LatestVersion() (*semver.Version, error) {
	file, err := g.Repository.GetFile(g.PassportPath())
	if err != nil {
		return nil, err
	}

	passport, err := Parse(file)
	if err != nil {
		return nil, err
	}

	return passport.Version, nil
}

func (g *Github) Validate() error {
	for _, t := range []*lo.Tuple2[string, string]{
		{"name", g.Repository.Name},
		{"owner", g.Repository.Owner},
		{"branch", g.Repository.Branch},
	} {
		if t.B == "" {
			return fmt.Errorf("missing required field in github.repo: %s", t.A)
		}
	}

	return nil
}
