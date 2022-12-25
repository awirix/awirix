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
   handler: function(query: string, progress: Progress): {Media}
   title: string|nil
   subtitle: string|nil
   placeholder: string|nil
   noun: Noun|nil
end

local record Layer
   handler: function(media: Media, progress: Progress): {Media}
   title: string|nil
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
   handler = function(query: string, progress: Progress): {Media}
      return {}
   end,
}

M.layers = {
	{
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
