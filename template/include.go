package template

import (
	"embed"
	_ "embed"
)

var (
	//go:embed include/editorconfig.tmpl
	templateEditorConfig []byte

	//go:embed include/extension.gitignore
	templateGitignore []byte
)

var (
	//go:embed include/lua
	luaTemplates embed.FS

	//go:embed include/fennel
	fennelTemplates embed.FS

	//go:embed include/typescript
	typescriptTemplates embed.FS
)
