package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/vivi-app/vivi/option"
	"net/http"
)

type Repository struct {
	Owner  string `toml:"owner"`
	Name   string `toml:"name"`
	Branch string `toml:"branch"`

	repo  option.Option[*github.Repository]
	files option.Option[[]*File]
}

func (r *Repository) URL() string {
	return fmt.Sprintf("https://github.com/%s/%s", r.Owner, r.Name)
}

func (r *Repository) SVNURL() (string, error) {
	if err := r.Setup(); err != nil {
		return "", err
	}

	return r.repo.MustGet().GetSVNURL(), nil
}

func (r *Repository) GetFile(path string) (*File, error) {
	files, err := r.Files()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.Path == path {
			return file, nil
		}
	}

	return nil, fmt.Errorf("file not found: %s", path)
}

func (r *Repository) Files() ([]*File, error) {
	if r.files.IsPresent() {
		return r.files.MustGet(), nil
	}

	err := r.Setup()
	if err != nil {
		return nil, err
	}

	tree, resp, err := client.Git.GetTree(context.Background(), r.Owner, r.Name, r.Branch, false)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	files := make([]*File, len(tree.Entries))
	for i, entry := range tree.Entries {
		files[i] = r.newFile(entry.GetPath(), entry.GetSHA())
	}

	r.files = option.Some(files)
	return files, nil
}

func (r *Repository) newFile(path, sha string) *File {
	return &File{
		Repository: r,
		Path:       path,
		SHA:        sha,
	}
}

func (r *Repository) Setup() error {
	if r.repo.IsPresent() {
		return nil
	}

	repo, resp, err := client.Repositories.Get(context.Background(), r.Owner, r.Name)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	r.repo = option.Some(repo)
	return nil
}
