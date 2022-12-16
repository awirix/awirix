package template

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/filename"
	"github.com/vivi-app/vivi/key"
)

type Preset int

const (
	PresetLua Preset = iota + 1
	PresetFennel
)

func (p Preset) String() string {
	switch p {
	case PresetLua:
		return "Lua"
	case PresetFennel:
		return "Fennel"
	default:
		return "Unknown"
	}
}

func PresetFromString(preset string) (Preset, bool) {
	switch preset {
	case PresetLua.String():
		return PresetLua, true
	case PresetFennel.String():
		return PresetFennel, true
	default:
		return 0, false
	}
}

func Generate(preset Preset) (templates map[string][]byte) {
	var tmpl = make(map[string][]byte)

	if viper.GetBool(key.ExtensionsTemplateEditorConfig) {
		tmpl[filename.EditorConfig] = EditorConfig()
	}

	switch preset {
	case PresetLua:
		tmpl[filename.Scraper] = LuaScraper()
		tmpl[filename.Tester] = LuaTester()

		if viper.GetBool(key.ExtensionsTemplateStylua) {
			tmpl[filename.Stylua] = Stylua()
		}
	case PresetFennel:
		const (
			scraperFnl = "scraper.fnl"
			testerFnl  = "tester.fnl"
		)

		runFennel := func(source string) []byte {
			return []byte(fmt.Sprintf(`return require('fennel').install().dofile('%s')`, source))
		}

		tmpl[scraperFnl] = FennelScraper()
		tmpl[filename.Scraper] = runFennel(scraperFnl)

		tmpl[testerFnl] = FennelTester()
		tmpl[filename.Tester] = runFennel(testerFnl)

		tmpl["fennel.lua"] = fennelSource
	}

	return tmpl
}
