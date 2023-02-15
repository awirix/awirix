package manager

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/extensions/passport"
	"github.com/awirix/awirix/filename"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/github"
	"github.com/awirix/awirix/text"
	"github.com/awirix/awirix/where"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

var (
	ErrInvalidURL            = errAdd(fmt.Errorf("invalid URL"))
	ErrRepoPassportMissing   = errAdd(fmt.Errorf("repository does not contain a %s", filename.PassportJSON))
	ErrRepoInvalidPassport   = errAdd(fmt.Errorf("repository does not contain a valid %s", filename.PassportJSON))
	ErrExtensionInstalled    = errAdd(fmt.Errorf("extension is already installed"))
	ErrRepoNamePrefixMissing = errAdd(fmt.Errorf("missing '%s' prefix in the repo name", app.Prefix))
)

var (
	githubURLRegex = regexp.MustCompile(`^https?://github.com/(?P<owner>[^/]+)/(?P<name>[^/]+)(?:\.git)?$`)
)

type AddOptions struct {
	URL          string
	SkipConfirm  bool
	SkipValidate bool
}

func errAdd(err error) error {
	return errors.Wrap(err, "add")
}

func confirm(msg string) (bool, error) {
	var confirm bool

	err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf(msg),
		Default: false,
	}, &confirm)

	return confirm, err
}

func Add(url string, options *AddOptions) (*extension.Extension, error) {
	if options == nil {
		options = &AddOptions{}
	}

	url = strings.TrimSuffix(url, ".git")
	if !githubURLRegex.MatchString(url) {
		return nil, ErrInvalidURL
	}

	groups := text.RegexpGroups(githubURLRegex, url)
	owner, name := groups["owner"], groups["name"]

	if !strings.HasPrefix(name, app.Prefix) {
		return nil, ErrRepoNamePrefixMissing
	}

	repo := github.Repository{Owner: owner, Name: name}
	err := repo.Setup()
	if err != nil {
		return nil, errAdd(err)
	}

	file, err := repo.GetFile(filename.PassportJSON)
	if err != nil {
		return nil, ErrRepoPassportMissing
	}

	data, err := file.Contents()
	if err != nil {
		return nil, err
	}

	thePassport, err := passport.New(bytes.NewBuffer(data))
	if err != nil {
		return nil, ErrRepoInvalidPassport
	}

	// TODO: add confirmation
	path := filepath.Join(where.Extensions(), thePassport.ID)

	exists, err := filesystem.Api().Exists(path)
	if err != nil {
		return nil, errAdd(err)
	}

	if exists {
		return nil, ErrExtensionInstalled
	}

	if exists {

	}

	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:          url,
		Depth:        1,
		SingleBranch: true,
	})

	if err != nil {
		return nil, errAdd(err)
	}

	return extension.New(path)
}
