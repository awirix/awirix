package templates

import (
	"bytes"
	"fmt"
	"io/fs"
	"text/template"
)

//go:generate enumer -type=Preset -trimprefix=Preset

type Preset int

const (
	// PresetLua is a preset for Lua
	PresetLua Preset = iota

	// PresetFennel is a preset for Fennel
	PresetFennel

	// PresetTypescript is a preset for Typescript
	PresetTypescript
)

type ErrInvalidPreset error

func Get(preset Preset, info Info) (map[string]*bytes.Buffer, error) {
	var f fs.FS

	switch preset {
	case PresetLua:
		f = FSLua
	case PresetFennel:
		f = FSFennel
	case PresetTypescript:
		f = FSTypescript
	default:
		return nil, ErrInvalidPreset(fmt.Errorf("invalid preset: %s", preset))
	}

	var tree = map[string]*bytes.Buffer{
		".gitignore":    bytes.NewBuffer(GitIgnore),
		".editorconfig": bytes.NewBuffer(EditorConfig),
	}

	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		b, err := fs.ReadFile(f, path)
		if err != nil {
			return err
		}

		t := template.New(preset.String())
		t, err = t.Parse(string(b))
		if err != nil {
			return err
		}

		var buf bytes.Buffer

		err = t.Execute(&buf, info)
		if err != nil {
			return err
		}

		tree[d.Name()] = &buf
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tree, nil
}
