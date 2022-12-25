package scraper

type Noun struct {
	singular,
	plural string
}

func (n *Noun) Singular() string {
	singular := n.singular
	if singular == "" {
		singular = "media"
	}

	return singular
}

func (n *Noun) Plural() string {
	plural := n.plural
	if plural == "" {
		plural = n.Singular() + "s"
	}

	return plural
}
