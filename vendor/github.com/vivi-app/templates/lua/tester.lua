{{- /*gotype:github.com/vivi-app/vivi/template.meta*/ -}}-- vim:ts=3 ss=3 sw=3 expandtab

local M = {}

--- Tests this extension.
function M.test()
	assert(2 + 2 == 4, 'Math is broken')
end

return M
