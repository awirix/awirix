package passport

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/github"
)

type Github struct {
	Repository github.Repository `toml:"repo"`
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
