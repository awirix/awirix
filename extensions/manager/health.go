package manager

import (
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/log"
	"github.com/awirix/awirix/where"
	"io"
	"path/filepath"
)

func CheckHealth(report io.Writer) {
	path := where.Extensions()

	dir, err := filesystem.Api().ReadDir(path)
	if err != nil {
		log.WriteErrorf(report, "failed to read extensions directory: %s", err)
		report.Write([]byte{0x0a})
		return
	}

	idMap := make(map[string]struct{})

	for _, owner := range dir {
		if !owner.IsDir() {
			continue
		}

		path := filepath.Join(path, owner.Name())
		dir, err := filesystem.Api().ReadDir(path)
		if err != nil {
			log.WriteErrorf(report, "failed to read extensions directory: %s", err)
			report.Write([]byte{0x0a})
			continue
		}

		for _, d := range dir {
			if !d.IsDir() {
				continue
			}

			extensionPath := filepath.Join(path, d.Name())
			ext, err := extension.New(extensionPath)
			if err != nil {
				log.WriteErrorf(report, "failed to load extension at '%s': %s", extensionPath, err)
				report.Write([]byte{0x0a})
				continue
			}

			if _, ok := idMap[ext.Passport().ID]; ok {
				log.WriteErrorf(report, "duplicate extension ID '%s' found", ext.Passport().ID)
				report.Write([]byte{0x0a})
				continue
			}

			idMap[ext.Passport().ID] = struct{}{}
		}
	}
}
