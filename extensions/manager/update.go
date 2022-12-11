package manager

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/key"
	"os"
	"time"
)

func UpdateExtension(ext *extension.Extension) (*extension.Extension, error) {
	theSpinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr), spinner.WithColor("cyan"))
	progress := func(text string) {
		theSpinner.Suffix = " " + text
	}

	progress(fmt.Sprintf("Updating %s", ext.String()))
	theSpinner.Start()
	defer theSpinner.Stop()

	path := ext.Path()

	// Will try to pull first, if that fails, it will clone
	if viper.GetBool(key.ExtensionsUpdateTryPull) {
		var repo *git.Repository
		repo, err := git.PlainOpen(path)

		if err == nil {
			err = updatePull(progress, ext, repo)
			if err == nil {
				return extension.New(ext.Path())
			}
		}
	}

	err := updateClone(progress, ext)
	if err != nil {
		return nil, err
	}

	updated, err := extension.New(ext.Path())
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func updatePull(progress func(string), ext *extension.Extension, repo *git.Repository) (err error) {
	if repo == nil {
		progress("Opening repository")
		repo, err = git.PlainOpen(ext.Path())
		if err != nil {
			return err
		}
	}

	passportRepo := ext.Passport().Repository
	pullOptions := &git.PullOptions{
		Depth: 1,
		Force: true,
	}

	if passportRepo != nil {
		pullOptions.RemoteURL = passportRepo.URL()
	}

	progress("Pulling changes")
	tree, err := repo.Worktree()
	if err != nil {
		return err
	}

	return tree.Pull(pullOptions)
}

func updateClone(progress func(string), ext *extension.Extension) error {
	if ext.Passport().Repository == nil {
		return fmt.Errorf("no repository specified in the passport")
	}

	path := ext.Path()

	tmpPath, err := filesystem.Api().TempDir("", ext.Passport().Name)
	if err != nil {
		return err
	}

	// ignore errors
	defer filesystem.Api().RemoveAll(tmpPath)

	progress("Cloning repository")
	_, err = git.PlainClone(tmpPath, false, &git.CloneOptions{
		URL:   ext.Passport().Repository.URL(),
		Depth: 1,
	})

	if err != nil {
		return err
	}

	progress("Moving files")
	_, err = extension.New(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to load passport: %w", err)
	}

	err = filesystem.Api().RemoveAll(path)
	if err != nil {
		return err
	}

	return filesystem.Api().Rename(tmpPath, path)
}
