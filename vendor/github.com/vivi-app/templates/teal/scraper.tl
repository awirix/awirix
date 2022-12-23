local record Media
	display: string
	about: string|nil
end

local type Progress = function(message: string)
local type Layer = function(media: Media, progress: Progress): {Media}

local record Module
	layers: {string:Layer}
	search: function(query: string, progress: Progress): {Media}
	prepare: function(media: Media, progress: Progress): Media
	stream: function(media: Media, progress: Progress)
	download: function(media: Media, progress: Progress)
end

local M: Module = {}

function M.search(query: string, progress: Progress): {Media}
   progress("Searching for " .. query)
	return {}
end

local function layer(media: Media, progress: Progress): {Media}
   progress("Layering " .. media.display)
	return {}
end

M.layers = {
	["First Layer"] = layer,
}

function M.prepare(media: Media, progress: Progress): Media
	progress("Preparing " .. media.display)
	return media
end

function M.stream(media: Media, progress: Progress)
	progress("Streaming " .. media.display)
end

function M.download(media: Media, progress: Progress)
	progress("Downloading " .. media.display)
end

return M
