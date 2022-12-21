-- vim:ts=3 ss=3 sw=3 expandtab

-- Learn Lua:   https://learnxinyminutes.com/docs/lua
-- Style Guide: https://github.com/luarocks/lua-style-guide

--- Table that have a string field named `{{ .Fields.Display }}` used for string representation
--- with optional `{{ .Fields.About }}` for brief description.
--- @alias media { {{ .Fields.Display }}: string, {{ .Fields.About }}: string|nil, [any]: any }

--- A function that is used to pass progress information to the vivi's ui.
--- @alias progress fun(message: string)

local M = {}

--- Searches for the media
--- This function may be omitted if this extension does not provide searching functionality.
--- Might be the case if it is dedicated to the single show/movie/book/...
--- @param query string The query to search for.
--- @param progress progress The progress writer
--- @return media[] # The media that match the query.
function M.search(query, progress)
   progress('Searching for ' .. query .. '...')
   return {}
end

--- Returns a list of nested medias for media entry.
--- E.g. it should return list of episodes for the show, chapters for the book and so on.
--- @param media media? The entry to explore.
--- @param progress progress The progress writer
--- @return media[] # The list of nested medias.
local function layer(media, progress)
   assert(media, 'media expected')
   progress('Layer ' .. media.display .. '...')
   return {}
end

--- Each layer returns a list of sub-media for the given one.
--- For example, you can search for a show, then selected show will be passed to the first layer that's responsible for returning show's seasons.
--- After that, the selected season will be passed to the second layer that would return season's episodes.
--- Layers may be omitted (nil or 0 length) if this extension does not provide such functionality (e.g. just search and watch, no seasons, no episodes).
--- If `search` function is omitted, first layer will receive a `nil` instead of the media.
--- @type table<string, fun(progress: progress, media: media|nil): media[]>[]
M.layers = {
   ['First Layer'] = layer
}

--- Prepares a media for streaming/downloading.
--- When you need to do some heavy operations before passing it further.
--- E.g. decode a video stream url, calculate a hash of the file, etc.
--- @param media media The media to prepare.
--- @param progress progress The progress writer
--- @return media # The prepared media
function M.prepare(media, progress)
   progress('Preparing ' .. media.display .. '...')
   return media
end

--- Streams a media
--- May be omitted if the extension does not provide streaming functionality.
--- However, at least one of the `{{ .Fn.Stream }}` or `{{ .Fn.Download }}` functions must be present.
--- @param media media The media to stream (watch, read, open).
--- @param progress progress The progress writer
function M.stream(media, progress)
   progress('Streaming ' .. media.display .. '...')
   error('Not implemented')

   -- Example: require('{{ .App }}').api.open(media.url)
end

--- Downloads a media.
--- May be omitted if this extension does not provide downloading functionality.
--- However, at least one of the `{{ .Fn.Stream }}` or `{{ .Fn.Download }}` functions must be present.
--- @param media media The media to download.
--- @param progress progress The progress writer
function M.download(media, progress)
   progress('Downloading ' .. media.display .. '...')
   error('Not implemented')

   -- Example: require('{{ .App }}').api.download(media.url, media.{{ .Fields.Display }})
end

return M
