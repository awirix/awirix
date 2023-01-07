package extension

import (
	"fmt"
	"github.com/awirix/awirix/extensions/passport"
	"github.com/awirix/awirix/filename"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/log"
	"github.com/awirix/awirix/scraper"
	"github.com/awirix/awirix/tester"
	"github.com/awirix/awirix/where"
	"github.com/awirix/lua"
	"path/filepath"
)

type Extension struct {
	path     string
	passport *passport.Passport
	scraper  *scraper.Scraper
	tester   *tester.Tester
	state    *lua.LState
}

func (e *Extension) loadPassport() error {
	path := filepath.Join(e.Path(), filename.Passport)

	exists, err := filesystem.Api().Exists(path)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("%s is missing", filename.Passport)
	}

	file, err := filesystem.Api().Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	thePassport, err := passport.New(file)
	if err != nil {
		return err
	}

	e.passport = thePassport
	return nil
}

func (e *Extension) LoadScraper(debug bool) error {
	e.initState(debug)

	file, err := filesystem.Api().Open(filepath.Join(e.Path(), filename.Scraper))
	if err != nil {
		return err
	}
	defer file.Close()

	theScraper, err := scraper.New(e.state, file)
	if err != nil {
		return err
	}

	theScraper.SetExtensionContext(&scraper.Context{
		Progress: func(message string) {
			log.Tracef("%s: progress: %s", e, message)
		},
		Error: func(err error) {
			log.Tracef("%s: error: %s", e, err)
		},
	})
	e.scraper = theScraper
	return nil
}

func (e *Extension) LoadTester(debug bool) error {
	e.initState(debug)

	file, err := filesystem.Api().Open(filepath.Join(e.Path(), filename.Tester))
	if err != nil {
		return err
	}
	defer file.Close()

	theTester, err := tester.New(e.state, file)
	if err != nil {
		return err
	}

	e.tester = theTester
	return nil
}

func (e *Extension) IsScraperLoaded() bool {
	return e.scraper != nil
}

func (e *Extension) String() string {
	return e.Passport().Name
}

func (e *Extension) Author() string {
	return filepath.Base(filepath.Dir(e.Path()))
}

func (e *Extension) Passport() *passport.Passport {
	return e.passport
}

func (e *Extension) Scraper() *scraper.Scraper {
	return e.scraper
}

func (e *Extension) Tester() *tester.Tester {
	return e.tester
}

func (e *Extension) Path() string {
	return e.path
}

func (e *Extension) Downloads() string {
	return filepath.Join(where.Downloads(), e.Passport().ID)
}
