local record Tester
	test: function()
end

local M: Tester = {}

function M.test()
	assert(2 + 2 == 4, 'Math is broken')
end

return M
