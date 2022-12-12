package template

import _ "embed"

//go:embed editorconfig.tmpl
var templateEditorConfig string

func EditorConfig() []byte {
	return []byte(templateEditorConfig)
}
