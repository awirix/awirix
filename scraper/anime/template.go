package anime

import "fmt"

const (
	fnSearch   = "__SEARCH"
	fnEpisodes = "__EPISODES"
	fnPrepare  = "__PREPARE"
	fnWatch    = "__WATCH"
	fnDownload = "__DOWNLOAD"
)

var Template = fmt.Sprintf(`--- Anime scraper
-- @script Scraper

--- Searches for anime.
-- @param query string The query to search for.
-- @return table of tables A list of anime. You can use any arbitrary keys, but the following are required:
--   - show: string The string to display in the search results.
function %[1]s(query)
	return {}
end

--- Returns a list of episodes for an anime.
-- @param anime table The anime to get episodes for.
-- @return table of tables A list of episodes. You can use any arbitrary keys, but the following are required:
--   - show: string The string to display in the episode list.
function %[2]s(anime)
	return {}
end

--- Prepares an episode for watching/downloading.
-- @param episode table The episode to prepare.
-- @return table The prepared episode. You can use any arbitrary keys
function %[3]s(episode)
	return episode
end

--- Watches an episode.
-- @param episode table The episode to watch.
function %[4]s(episode)
	require("vivi").watch(episode.url)
end

--- Downloads an episode.
-- @param episode table The episode to download.
function %[5]s(episode)
	require("vivi").download(episode.show, episode.url)
end
`, fnSearch, fnEpisodes, fnPrepare, fnWatch, fnDownload)
