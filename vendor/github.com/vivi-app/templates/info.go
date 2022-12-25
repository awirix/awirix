package templates

type Info struct {
	// Name is the name of the scraper
	Name string

	// About is a short description of the scraper
	About string

	// NSFW is a flag indicating whether the scraper is NSFW
	NSFW bool

	// Libs is the paths to the libraries used by the lua scraper
	Libs []string
}
