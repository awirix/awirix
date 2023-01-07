package github

import (
	"context"
	"fmt"
	"github.com/awirix/awirix/option"
	"net/http"
)

type File struct {
	Repository *Repository
	Path       string
	SHA        string
	contents   *option.Option[[]byte]
}

func (f *File) URL() string {
	return fmt.Sprintf("%s/blob/%s/%s", f.Repository.URL(), f.Repository.Branch, f.Path)
}

func (f *File) Setup() error {
	if f.contents.IsPresent() {
		return nil
	}

	err := f.Repository.Setup()
	if err != nil {
		return err
	}

	contents, resp, err := client.Git.GetBlobRaw(context.Background(), f.Repository.Owner, f.Repository.Name, f.SHA)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	f.contents = option.Some(contents)
	return nil
}

func (f *File) Contents() ([]byte, error) {
	err := f.Setup()
	if err != nil {
		return nil, err
	}

	return f.contents.MustGet(), nil
}
