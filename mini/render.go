package mini

import (
	"github.com/awirix/awirix/core"
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/style"
)

func renderMedia(media *core.Media) (rendered string) {
	rendered += media.String()

	if description := media.Description(); description != "" {
		rendered += " " + style.Faint(description)
	}

	return
}

func renderExtension(ext *extension.Extension) (rendered string) {
	rendered += ext.String()

	if about := ext.Passport().About; about != "" {
		rendered += " " + style.Faint(about)
	}

	return
}
