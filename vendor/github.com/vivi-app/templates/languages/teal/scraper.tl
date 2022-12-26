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
   title: string
   noun: Noun|nil
end

local record Action
   handler: function(media: Media, progress: Progress)
   title: string
   description: string|nil
end

local record Module
   search: Search|nil
	layers: {Layer}|nil
   actions: {Action}|nil
end

local M: Module = {}

M.search = {
   handler = function(query: string, progress: Progress): {Media}
      return {}
   end,
}

M.layers = {
	{
      title = 'Layer',
      handler = function(media: Media, progress: Progress): {Media}
         return {}
      end,
   }
}

M.actions = {
   {
      title = 'Stream',
      handler = function(media: Media, progress: Progress)
         error('Not implemented')
      end,
   },
   {
      title = 'Download',
      handler = function(media: Media, progress: Progress)
         error('Not implemented')
      end,
   },
}

return M
