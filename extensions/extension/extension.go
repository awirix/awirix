package extension

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/awirix/awirix/core"
	"github.com/awirix/awirix/extensions/passport"
	"github.com/awirix/awirix/filename"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/key"
	"github.com/awirix/awirix/log"
	"github.com/awirix/awirix/tester"
	"github.com/awirix/awirix/where"
	"github.com/awirix/lua"
	"github.com/enescakir/emoji"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func errExtension(err error) error {
	return errors.Wrap(err, "extension")
}

func errMissing(name string) error {
	return errExtension(fmt.Errorf("%s is missing", name))
}

type Extension struct {
	path     string
	passport *passport.Passport
	core     *core.Core
	tester   *tester.Tester
	state    *lua.LState
}

func (e *Extension) loadPassport() error {
	path := filepath.Join(e.Path(), filename.PassportJSON)

	exists, err := filesystem.Api().Exists(path)
	if err != nil {
		return errExtension(err)
	}

	if !exists {
		return errMissing(filename.PassportJSON)
	}

	file, err := filesystem.Api().Open(path)
	if err != nil {
		return errExtension(err)
	}
	defer file.Close()

	thePassport, err := passport.New(file)
	if err != nil {
		return errExtension(err)
	}

	e.passport = thePassport
	return nil
}

func (e *Extension) LoadCore(debug bool) error {
	e.initState(debug)

	path := filepath.Join(e.Path(), filename.MainLua)
	exists, err := filesystem.Api().Exists(path)
	if err != nil {
		return errExtension(err)
	}

	if !exists {
		return errMissing(filename.MainLua)
	}

	file, err := filesystem.Api().Open(path)
	if err != nil {
		return errExtension(err)
	}
	defer file.Close()

	theCore, err := core.New(e.state, file)
	if err != nil {
		return errExtension(err)
	}

	theCore.SetExtensionContext(&core.Context{
		Progress: func(message string) {
			log.Tracef("%s: progress: %s", e, message)
		},
		Error: func(err error) {
			log.Tracef("%s: error: %s", e, err)
		},
	})
	e.core = theCore
	return nil
}

func (e *Extension) LoadTester(debug bool) error {
	e.initState(debug)

	path := filepath.Join(e.Path(), filename.TestLua)
	exists, err := filesystem.Api().Exists(path)
	if err != nil {
		return errExtension(err)
	}

	if !exists {
		return errMissing(filename.TestLua)
	}

	file, err := filesystem.Api().Open(path)
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

func (e *Extension) IsCoreLoaded() bool {
	return e.core != nil
}

func (e *Extension) String() string {
	var b strings.Builder

	b.WriteString(e.Passport().Name)

	if viper.GetBool(key.IconShowExtensionFlag) {
		flag, err := emoji.CountryFlag(e.Passport().Language.Code)
		if err == nil {
			b.WriteRune(' ')
			b.WriteString(flag.String())
		}
	}

	return b.String()

	// having icon as a prefix looks better, but it causes problems with tui list filtering
	//return e.Passport().Icon.String() + " " + e.Passport().Name
}

func (e *Extension) Passport() *passport.Passport {
	return e.passport
}

func (e *Extension) Core() *core.Core {
	return e.core
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

func (e *Extension) Cache() string {
	path := filepath.Join(e.Path(), ".cache")
	lo.Must0(filesystem.Api().MkdirAll(path, os.ModePerm))
	return path
}
