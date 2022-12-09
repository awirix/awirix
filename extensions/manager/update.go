package manager

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/filesystem"
	"os"
	"time"
)

func UpdateExtension(ext *extension.Extension) error {
	theSpinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr), spinner.WithColor("cyan"))
	progress := func(text string) {
		theSpinner.Suffix = " " + text
	}

	progress(fmt.Sprintf("Updating %s", ext.String()))
	theSpinner.Start()
	defer theSpinner.Stop()

	if err := ext.LoadPassport(); err != nil {
		return err
	}

	path := ext.Path()
	repo, err := git.PlainOpen(path)

	if err == nil {
		err = updatePull(progress, ext, repo)
		if err == nil {
			return nil
		}
	}

	return updateClone(progress, ext)
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

	// if pull failed, try to remove and download again
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
	cloned := extension.New(tmpPath)
	cloned.Init()
	err = cloned.LoadPassport()
	if err != nil {
		return fmt.Errorf("failed to load passport: %w", err)
	}

	err = filesystem.Api().RemoveAll(path)
	if err != nil {
		return err
	}

	return filesystem.Api().Rename(tmpPath, path)
}
