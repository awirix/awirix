local record Context
   progress: function(message: string)
   error: function(message: string)
end

local record Noun
   singular: string
   plural: string|nil
end

local record Media
	title: string
	description: string|nil
end

local record Search
   handler: function(query: string, ctx: Context): {Media}
   title: string|nil
   subtitle: string|nil
   placeholder: string|nil
   noun: Noun|nil
end

local record Layer
   handler: function(media: Media, ctx: Context): {Media}
   title: string
   noun: Noun|nil
end

local record Action
   handler: function(media: {Media}, ctx: Context)
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
   handler = function(query: string, ctx: Context): {Media}
      return {}
   end,
}

M.layers = {
	{
      title = 'Layer',
      handler = function(media: Media, ctx: Context): {Media}
         return {}
      end,
   }
}

M.actions = {
   {
      title = 'Stream',
      handler = function(media: {Media}, ctx: Context)
         ctx.error('Not implemented')
      end,
   },
   {
      title = 'Download',
      multiple = true,
      handler = function(media: {Media}, ctx: Context)
         ctx.error('Not implemented')
      end,
   },
}

return M
