package extension

import (
	"context"
	"fmt"
	lua "github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/extensions/passport"
	"github.com/vivi-app/vivi/filename"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/tester"
	"github.com/vivi-app/vivi/vm"
	"path/filepath"
	"strings"
)

type Extension struct {
	path     string
	passport *passport.Passport
	scraper  *scraper.Scraper
	tester   *tester.Tester
	state    *lua.LState
}

func (e *Extension) initState(options *vm.Options) {
	if e.state != nil {
		return
	}

	state := vm.New(options)

	// this is hideous, but it works
	state.SetContext(context.WithValue(context.Background(), true, e))

	// add local files to the path
	pkg := state.GetGlobal("package").(*lua.LTable)
	paths := strings.Split(pkg.RawGetString("path").String(), ";")
	viviPaths := []string{
		filepath.Join(e.Path(), "?.lua"),
	}
	paths = append(viviPaths, paths...)
	pkg.RawSetString("path", lua.LString(strings.Join(paths, ";")))

	e.state = state
}

func (e *Extension) loadPassport() error {
	file, err := filesystem.Api().Open(filepath.Join(e.Path(), filename.Passport))
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

func (e *Extension) LoadScraper(options *vm.Options) error {
	e.initState(options)

	file, err := filesystem.Api().Open(filepath.Join(e.Path(), filename.Scraper))
	if err != nil {
		return err
	}
	defer file.Close()

	theScraper, err := scraper.New(context.Background(), e.state, file)
	if err != nil {
		return err
	}

	theScraper.SetProgress(func(string) {})
	e.scraper = theScraper
	return nil
}

func (e *Extension) LoadTester() error {
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

func (e *Extension) String() string {
	var name string
	if e.Passport() != nil {
		name = e.Passport().Name
	} else {
		name = filepath.Base(e.Path())
	}

	return fmt.Sprintf("%s/%s", e.Author(), name)
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
