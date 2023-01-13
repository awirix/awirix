package tui

import (
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/scraper"
)

type (
	msgExtensionLoaded  *extension.Extension
	msgSearchDone       []*scraper.Media
	msgLayerDone        []*scraper.Media
	msgLayerItemsSet    struct{}
	msgActionDone       *scraper.Action
	msgMediaInfoDone    struct{}
	msgError            error
	msgExtensionRemoved *extension.Extension
)
