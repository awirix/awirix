local type Progress = function(message: string)

local record Noun
   singular: string
   plural: string|nil
end

local record Media
	title: string
	description: string|nil
end

local record Search
   title: string
   handler: function(query: string, progress: Progress): {Media}
   noun: Noun|nil
end

local record Layer
   title: string
   handler: function(media: Media, progress: Progress): {Media}
   noun: Noun|nil
end

local record Module
   search: Search
	layers: {Layer}

	prepare: function(media: Media, progress: Progress): Media
	stream: function(media: Media, progress: Progress)
	download: function(media: Media, progress: Progress)
end

local M: Module = {}

M.search = {
   title = "Search",
   handler = function(query: string, progress: Progress): {Media}
      return {}
   end,
}

M.layers = {
	{
      title = "First Layer",
      handler = function(media: Media, progress: Progress): {Media}
         return {}
      end,
   }
}

function M.prepare(media: Media, progress: Progress): Media
	return media
end

function M.stream(media: Media, progress: Progress)
end

function M.download(media: Media, progress: Progress)
end

return M
