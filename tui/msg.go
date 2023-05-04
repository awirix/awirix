package tui

import (
	"github.com/awirix/awirix/core"
	"github.com/awirix/awirix/extensions/extension"
)

type (
	msgExtensionLoaded  *extension.Extension
	msgSearchDone       []*core.Media
	msgLayerDone        []*core.Media
	msgLayerItemsSet    struct{}
	msgActionDone       *core.Action
	msgMediaInfoDone    struct{}
	msgError            error
	msgExtensionRemoved *extension.Extension
)
