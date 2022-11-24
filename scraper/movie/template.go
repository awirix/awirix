package movie

import "fmt"

const (
	fnSearch   = "__SEARCH"
	fnPrepare  = "__PREPARE"
	fnWatch    = "__WATCH"
	fnDownload = "__DOWNLOAD"
)

// goland:language=lua
var Template = fmt.Sprintf(`
function %[1]s(query)
	return {}
end

function %[2]s(movie)
	return movie
end

function %[3]s(movie)
	require("vivi").watch(movie.url)
end

function %[4]s(movie)
    require("vivi").download(movie.url)
end
`, fnSearch, fnPrepare, fnWatch, fnDownload)
