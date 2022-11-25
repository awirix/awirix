package test

import (
	"fmt"
	"github.com/vivi-app/vivi/constant"
)

var Template = fmt.Sprintf(`local %[1]s = {}
local %[3]s = require('%[3]s')

function %[1]s.%[2]s()
	return nil
end

return %[1]s`,
	constant.ModuleTest,
	constant.FunctionTest,
	constant.ModuleScraper,
)
