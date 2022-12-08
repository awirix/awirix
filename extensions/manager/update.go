package manager

import (
	"github.com/go-git/go-git/v5"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/filesystem"
)

func UpdateExtension(ext *extension.Extension) error {
	if err := ext.LoadPassport(); err != nil {
		return err
	}

	path := ext.Path()
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	passportRepo := ext.Passport().Repository
	pullOptions := &git.PullOptions{
		Depth: 1,
		Force: true,
	}

	if passportRepo != nil {
		pullOptions.RemoteURL = passportRepo.URL()
	}

	tree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = tree.Pull(pullOptions)
	if err == nil {
		return nil
	}

	if passportRepo == nil {
		return err
	}

	// if pull failed, try to remove and download again
	tmpPath := path + ".tmp"

	// ignore errors
	defer filesystem.Api().RemoveAll(tmpPath)

	_, err = git.PlainClone(tmpPath, false, &git.CloneOptions{
		URL:   ext.Passport().Repository.URL(),
		Depth: 1,
	})

	if err != nil {
		return err
	}

	cloned := extension.New(tmpPath)
	cloned.Init()
	err = cloned.LoadPassport()
	if err != nil {
		return err
	}

	err = filesystem.Api().RemoveAll(path)
	if err != nil {
		return err
	}

	err = filesystem.Api().Rename(tmpPath, path)
	if err != nil {
		return err
	}

	return nil
}
