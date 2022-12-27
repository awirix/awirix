--- Name of the media type to display
--- @alias Noun { singular: string?, plural: string? }

--- Table that have a string field named `title` used for string representation
--- with optional `description` for a brief description.
--- @alias Media { title: string, description: string?, [any]: any }

--- Searches for the media
--- @alias Search { title: string?, subtitle: string?, placholder: string?, handler: fun(query: string, ctx: Context): Media[], noun: Noun? }

--- Each layer returns a list of sub-media for the given one.
--- For example, you can search for a show, then selected show will be passed to the first layer that's responsible for returning show's seasons.
--- After that, the selected season will be passed to the second layer that would return season's episodes.
--- @alias Layer { title: string, handler: fun(media: Media?, ctx: Context): Media[], noun: Noun? }[]


--- Actions are further actions that can be performed on the selected media.
--- Something like 'Stream' or 'Download'
--- @alias Action { title: string, handler: fun(media: Media[], ctx: Context), description: string? }

--- Context that is passed to the handler functions to report progress and errors.
--- @alias Context { progress: fun(message: string), error: fun(message: string) }

local M = {}

--- This step may be omitted if this extension does not provide searching functionality.
--- Might be the case if it is dedicated to the single show/movie/book/...
--- @type Search
M.search = {
   handler = function(query, ctx) return {} end
}

--- Layers may be omitted (nil or 0 length) if this extension does not provide such functionality (e.g. just search and watch, no seasons, no episodes).
--- If `search` function is omitted, first layer will receive a `nil` instead of the media.
--- @type Layer[]
M.layers = {
   {
      title = 'Layer',
      handler = function(media, ctx) return {} end
   }
}

--- Actions may be omitted (nil or 0 length) if this extension does not provide such functionality (e.g. just media browsing, no actions).
--- @type Action[]
M.actions = {
   {
      title = 'Search',
      handler = function (medias, ctx)
         ctx.error('Not implemented')
      end
   },
   {
      title = 'Download',
      handler = function (medias, ctx)
         ctx.error('Not implemented')
      end
   }
}


return M
