package template

import (
	"embed"
	"fmt"
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/filename"
	"github.com/vivi-app/vivi/key"
	"io/fs"
)

type Preset int

const (
	PresetLua Preset = iota + 1
	PresetFennel
	PresetTypescript
)

func (p Preset) String() string {
	switch p {
	case PresetLua:
		return "Lua"
	case PresetFennel:
		return "Fennel"
	case PresetTypescript:
		return "Typescript"
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
	case PresetTypescript.String():
		return PresetTypescript, true
	default:
		return 0, false
	}
}

func Generate(preset Preset) (map[string][]byte, error) {
	var tmpl = make(map[string][]byte)

	if viper.GetBool(key.ExtensionsTemplateEditorConfig) {
		tmpl[filename.EditorConfig] = templateEditorConfig
	}

	if viper.GetBool(key.ExtensionsTemplateGitignore) {
		tmpl[filename.Gitignore] = templateGitignore
	}

	bind := func(f *embed.FS, m map[string][]byte) error {
		return fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			contents, err := f.ReadFile(path)
			if err != nil {
				return err
			}

			m[d.Name()] = execTemplate(string(contents))

			return nil
		})
	}

	switch preset {
	case PresetLua:
		err := bind(&luaTemplates, tmpl)
		if err != nil {
			return nil, err
		}
	case PresetFennel:
		err := bind(&fennelTemplates, tmpl)
		if err != nil {
			return nil, err
		}
	case PresetTypescript:
		err := bind(&typescriptTemplates, tmpl)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown preset")
	}

	return tmpl, nil
}
