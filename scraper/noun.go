package scraper

type Noun struct {
	noun `lua:",squash"`
}

type noun struct {
	Singular,
	Plural string
}

func (n *Noun) Singular() string {
	if n.noun.Singular != "" {
		return n.noun.Singular
	}

	return "media"
}

func (n *Noun) Plural() string {
	if n.noun.Plural != "" {
		return n.noun.Plural
	}

	return n.Singular() + "s"
}
