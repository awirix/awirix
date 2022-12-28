package tui

import (
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/scraper"
)

type (
	msgExtensionLoaded *extension.Extension
	msgSearchDone      []*scraper.Media
	msgLayerDone       []*scraper.Media
	msgLayerItemsSet   struct{}
	msgActionDone      *scraper.Action
	msgMediaInfoDone   string
	msgError           error
)
