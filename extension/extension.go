package extension

import (
	"fmt"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/passport"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/util"
	"github.com/vivi-app/vivi/where"
	"path/filepath"
)

type Extension struct {
	passport *passport.Passport
	scraper  *scraper.Scraper
}

func (e *Extension) LoadScraper() error {
	theScraper, err := scraper.FromPath(e.Path())
	if err != nil {
		return err
	}

	e.scraper = theScraper
	return nil
}

func (e *Extension) String() string {
	return e.Passport().Name
}

func (e *Extension) Passport() *passport.Passport {
	return e.passport
}

func (e *Extension) Scraper() *scraper.Scraper {
	return e.scraper
}

func (e *Extension) IsInstalled() bool {
	path := e.Path()
	installed, err := filesystem.Api().ReadDir(where.Extensions())
	if err != nil {
		return false
	}

	for _, file := range installed {
		if filepath.Join(where.Extensions(), file.Name()) == path {
			return true
		}
	}

	return false
}

func (e *Extension) Path() string {
	dir := util.SanitizeFilename(e.Passport().ID)
	return filepath.Join(where.Extensions(), dir)
}

func (e *Extension) Install() error {
	svn, err := e.Passport().Github.Repository.SVNURL()
	if err != nil {
		return err
	}

	fmt.Println(svn)

	return nil
}

func (e *Extension) Uninstall() error {
	return filesystem.Api().RemoveAll(e.Path())
}
