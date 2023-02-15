package manager

import (
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/log"
	"github.com/awirix/awirix/where"
	"io"
	"path/filepath"
)

func Health(report io.Writer) {
	var newline = []byte{0x0a}
	path := where.Extensions()

	dir, err := filesystem.Api().ReadDir(path)
	if err != nil {
		log.WriteErrorf(report, "failed to read extensions directory: %s", err)
		report.Write(newline)
		return
	}

	idMap := make(map[string]struct{})

	var errOccurred bool
	for _, f := range dir {
		if !f.IsDir() {
			continue
		}

		path := filepath.Join(path, f.Name())
		ext, err := extension.New(path)
		if err != nil {
			errOccurred = true
			log.WriteErrorf(report, "failed to load extension at %q: %s", path, err)
			report.Write(newline)
			continue
		}

		id := ext.Passport().ID
		if _, ok := idMap[id]; ok {
			errOccurred = true
			log.WriteErrorf(report, "extension %q is installed multiple times", id)
			report.Write(newline)
		}

		idMap[id] = struct{}{}
	}

	if !errOccurred {
		log.WriteSuccessf(report, "Everything is OK")
		report.Write(newline)
	}
}
