# Vivi

> One scraper to rule them all!

Vivi is currently a **heavy work in progress** project that aims to replace **all** existing media scrapers.
Watch movies, anime, read manga, download torrents, you name it!
Literally limitless possibilites due to *extension-based* structure powered by *Lua Scripts*.

At this stage, there's not much to show (in fact, by "not much" I mean nothing),
it's not even an alpha version, more like a draft, so you probably don't want to use it right now.

But you can take a look at the scraper's template to get a brief idea what is going on

```lua
--- Table that have a string field named `display` used for string representation
--- with optional `about` for brief description.
--- @alias media { display: string, about: string|nil, [any]: any }

--- A function that is used to pass progress information to the vivi's ui.
--- @alias progress fun(message: string)

local M = {}

--- Searches for the media
--- This function may be omitted if this extension does not provide searching functionality.
--- Might be the case if it is dedicated to the single show/movie/book/...
--- @param progress progress @The progress writer
--- @param query string @The query to search for.
--- @return media[] # The media that match the query.
function M.search(progress, query)
   progress('Searching for ' .. query .. '...')
   return {}
end

--- Returns a list of nested medias for media entry.
--- E.g. it should return list of episodes for the show, chapters for the book and so on.
--- @param progress progress @The progress writer
--- @param media media? @The entry to explore.
--- @return media[] # The list of nested medias.
local function layer(progress, media)
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
   ["First Layer"] = layer
}

--- Prepares a media for streaming/downloading.
--- When you need to do some heavy operations before passing it further.
--- E.g. decode a video stream url, calculate a hash of the file, etc.
--- @param progress progress @The progress writer
--- @param media media @The media to prepare.
--- @return media # The prepared media
function M.prepare(progress, media)
   progress('Preparing ' .. media.display .. '...')
   return media
end

--- Streams a media
--- May be omitted if the extension does not provide streaming functionality.
--- However, at least one of the `stream` or `download` functions must be present.
--- @param progress progress @The progress writer
--- @param media media @The media to stream (watch, read, open).
function M.stream(progress, media)
   progress('Streaming ' .. media.display .. '...')
   error('Not implemented')

   -- Example: require('vivi').api.open(media.url)
end

--- Downloads a media.
--- May be omitted if this extension does not provide downloading functionality.
--- However, at least one of the `stream` or `download` functions must be present.
--- @param progress progress @The progress writer
--- @param media media @The media to download.
function M.download(progress, media)
   progress('Downloading ' .. media.display .. '...')
   error('Not implemented')

   -- Example: require('vivi').api.download(media.url, media.display)
end

return M
```
