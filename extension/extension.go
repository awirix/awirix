package extension

import (
	"context"
	"fmt"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/lualib"
	"github.com/vivi-app/vivi/passport"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/tester"
	"github.com/vivi-app/vivi/util"
	"github.com/vivi-app/vivi/where"
	lua "github.com/yuin/gopher-lua"
	"path/filepath"
	"strings"
)

type Extension struct {
	passport *passport.Passport
	scraper  *scraper.Scraper
	tester   *tester.Tester
	state    *lua.LState
}

func (e *Extension) Init() {
	if e.state != nil {
		return
	}

	L := lua.NewState()
	L.SetContext(context.WithValue(context.Background(), "passport", e.Passport()))

	// Load the standard libraries except for the debug, io and os
	for _, openLib := range []lua.LGFunction{
		lua.OpenBase,
		lua.OpenTable,
		lua.OpenString,
		lua.OpenMath,
		lua.OpenCoroutine,
		lua.OpenChannel,
		lua.OpenPackage,
	} {
		openLib(L)
	}

	pkg := L.GetGlobal("package").(*lua.LTable)
	paths := strings.Split(pkg.RawGetString("path").String(), ";")

	viviPaths := []string{
		filepath.Join(e.Path(), "?.lua"),
	}

	paths = append(viviPaths, paths...)

	pkg.RawSetString("path", lua.LString(strings.Join(paths, ";")))

	lualib.Preload(L)
	e.state = L
}

func (e *Extension) LoadScraper() error {
	file, err := filesystem.Api().Open(filepath.Join(e.Path(), constant.Scraper))
	if err != nil {
		return err
	}
	defer file.Close()

	theScraper, err := scraper.New(e.state, file)
	if err != nil {
		return err
	}

	e.scraper = theScraper
	return nil
}

func (e *Extension) LoadTester() error {
	file, err := filesystem.Api().Open(filepath.Join(e.Path(), constant.Tester))
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
	return e.Passport().Name
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
