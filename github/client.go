package github

import (
	"github.com/google/go-github/v48/github"
)

var client *github.Client

func init() {
	client = github.NewClient(nil)
}
