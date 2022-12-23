package templates

import "embed"

var (
	//go:embed .gitignore
	gitIgnore []byte

	//go:embed .editorconfig
	editorConfig []byte
)

var (
	//go:embed lua
	fsLua embed.FS

	//go:embed fennel
	fsFennel embed.FS

	//go:embed typescript
	fsTypescript embed.FS

	//go:embed teal
	fsTeal embed.FS

	//go:embed yue
	fsYue embed.FS
)
