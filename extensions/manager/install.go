package manager

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/extensions/passport"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/github"
	"github.com/vivi-app/vivi/util"
	"github.com/vivi-app/vivi/where"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type InstallOptions struct {
	Progress     io.Writer
	URL          string
	SkipConfirm  bool
	SkipValidate bool
}

func InstallExtension(options *InstallOptions) (*extension.Extension, error) {
	if !util.IsURL(options.URL) {
		return nil, fmt.Errorf("invalid URL")
	}

	trimmed := strings.TrimSuffix(options.URL, ".git")
	repoName := filepath.Base(trimmed)
	repoOwner := filepath.Base(filepath.Dir(trimmed))

	path := filepath.Join(where.Extensions(), util.SanitizeFilename(repoOwner), util.SanitizeFilename(repoName))

	if exists, err := filesystem.Api().Exists(path); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("extension already installed: %s/%s", repoOwner, repoName)
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	progress := func(text string) {
		s.Suffix = " " + text
	}
	progress(" Preparing...")
	s.Start()
	defer s.Stop()

	if !options.SkipValidate {
		repo := github.Repository{
			Owner: repoOwner,
			Name:  repoName,
		}

		progress("Getting repository information...")

		err := repo.Setup()
		if err != nil {
			return nil, err
		}

		progress("Searching for " + constant.FilenamePassport)
		file, err := repo.GetFile(constant.FilenamePassport)
		if err != nil {
			return nil, fmt.Errorf("repository does not contain a %s", constant.FilenamePassport)
		}

		progress("Reading " + constant.FilenamePassport)
		data, err := file.Contents()
		if err != nil {
			return nil, err
		}

		progress("Parsing " + constant.FilenamePassport)
		thePassport, err := passport.New(bytes.NewBuffer(data))
		if err != nil {
			return nil, fmt.Errorf("repository does not contain a valid passport: %s", err)
		}

		if !options.SkipConfirm {
			s.Stop()

			fmt.Println(thePassport.Info())
			fmt.Println()

			var confirm bool

			err := survey.AskOne(&survey.Confirm{
				Message: fmt.Sprintf("Install?"),
				Default: false,
			}, &confirm)

			if err != nil {
				return nil, err
			}

			if !confirm {
				return nil, fmt.Errorf("installation cancelled")
			}
		}

		s.Start()
	}

	progress("Cloning repository...")
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      options.URL,
		Progress: options.Progress,
		Depth:    1,
	})

	if err != nil {
		return nil, err
	}

	s.Stop()

	ext := extension.New(path)

	if !options.SkipValidate {
		err = ext.LoadPassport()
		if err != nil {
			return nil, err
		}
	}

	return ext, nil
}
