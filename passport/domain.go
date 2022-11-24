package passport

type Domain string

const (
	DomainAnime  = Domain("anime")
	DomainMovies = Domain("movies")
	DomainShows  = Domain("shows")
)

var Domains = []Domain{
	DomainAnime,
	DomainMovies,
	DomainShows,
}
