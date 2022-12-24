-- vim:ts=3 ss=3 sw=3 expandtab

-- Learn Lua:   https://learnxinyminutes.com/docs/lua
-- Style Guide: https://github.com/luarocks/lua-style-guide

--- Name of the media type to display
--- @alias noun { singular: string, plural: string? }

--- Table that have a string field named `title` used for string representation
--- with optional `description` for a brief description.
--- @alias media { title: string, description: string?, [any]: any }

--- A function that is used to pass progress information to the vivi's ui.
--- @alias progress fun(message: string)

local M = {}

--- Searches for the media
--- This step may be omitted if this extension does not provide searching functionality.
--- Might be the case if it is dedicated to the single show/movie/book/...
--- @type { title: string, handler: fun(query: string, progress: progress): media[], noun: noun? }
M.search = {
   title = "Search",
   handler = function(query, progress) return {} end
}

--- Each layer returns a list of sub-media for the given one.
--- For example, you can search for a show, then selected show will be passed to the first layer that's responsible for returning show's seasons.
--- After that, the selected season will be passed to the second layer that would return season's episodes.
--- Layers may be omitted (nil or 0 length) if this extension does not provide such functionality (e.g. just search and watch, no seasons, no episodes).
--- If `search` function is omitted, first layer will receive a `nil` instead of the media.
--- @type { title: string, handler: fun(media: media?, progress: progress): media[], noun: noun? }[]
M.layers = {
   {
      title = "First Layer",
      handler = function(media, progress) return {} end
   }
}

--- Prepares a media for streaming/downloading.
--- When you need to do some heavy operations before passing it further.
--- E.g. decode a video stream url, calculate a hash of the file, etc.
--- @param media media The media to prepare.
--- @param progress progress The progress writer
--- @return media # The prepared media
function M.prepare(media, progress)
   return media
end

--- Streams a media
--- May be omitted if the extension does not provide streaming functionality.
--- However, at least one of the `stream` or `download` functions must be present.
--- @param media media The media to stream (watch, read, open).
--- @param progress progress The progress writer
function M.stream(media, progress)
   error('Not implemented')
end

--- Downloads a media.
--- May be omitted if this extension does not provide downloading functionality.
--- However, at least one of the `stream` or `download` functions must be present.
--- @param media media The media to download.
--- @param progress progress The progress writer
function M.download(media, progress)
   error('Not implemented')
end

return M
