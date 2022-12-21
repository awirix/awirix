package templates

import "embed"

var (
	//go:embed .gitignore
	GitIgnore []byte

	//go:embed .editorconfig
	EditorConfig []byte
)

var (
	//go:embed lua
	FSLua embed.FS

	//go:embed fennel
	FSFennel embed.FS

	//go:embed typescript
	FSTypescript embed.FS
)
