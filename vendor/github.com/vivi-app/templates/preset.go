package templates

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

//go:generate enumer -type=Preset -trimprefix=Preset

type Preset int

const (
	// PresetLua is a preset for Lua language
	PresetLua Preset = iota + 1

	// PresetFennel is a preset for Fennel language
	PresetFennel

	// PresetTypescript is a preset for Typescript language
	PresetTypescript

	// PresetTeal is a preset for the Teal language
	PresetTeal

	// PresetYue is a preset fot the Yue language
	PresetYue
)

var filesystems = map[Preset]fs.FS{
	PresetLua:        fsLua,
	PresetFennel:     fsFennel,
	PresetTypescript: fsTypescript,
	PresetTeal:       fsTeal,
	PresetYue:        fsYue,
}

// ErrInvalidPreset is returned when an invalid preset is given
type ErrInvalidPreset error

// Get gets the template tree for the given preset.
// Returned map is a map of file names to file contents
func Get(preset Preset, info Info) (map[string]*bytes.Buffer, error) {
	f, ok := filesystems[preset]
	if !ok {
		return nil, ErrInvalidPreset(fmt.Errorf("invalid preset: %s", preset))
	}

	var tree = map[string]*bytes.Buffer{
		".gitignore":    bytes.NewBuffer(gitIgnore),
		".editorconfig": bytes.NewBuffer(editorConfig),
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

		var (
			buf  bytes.Buffer
			name string
		)

		// Execute template only if it's a .tmpl file
		if filepath.Ext(d.Name()) == ".tmpl" {
			name = strings.TrimSuffix(d.Name(), ".tmpl")

			t := template.New(preset.String())
			t, err = t.Parse(string(b))
			if err != nil {
				return err
			}

			err = t.Execute(&buf, info)
			if err != nil {
				return err
			}
		} else {
			name = d.Name()
			buf.Write(b)
		}

		tree[name] = &buf
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tree, nil
}
