package templates

import "embed"

var (
	//go:embed .gitignore
	gitIgnore []byte

	//go:embed .editorconfig
	editorConfig []byte
)

var (
	//go:embed languages/lua
	fsLua embed.FS

	//go:embed languages/fennel
	fsFennel embed.FS

	//go:embed languages/typescript
	fsTypescript embed.FS

	//go:embed languages/teal
	fsTeal embed.FS

	//go:embed languages/yue
	fsYue embed.FS
)
