package scraper

import "io"

type Scraper struct{}

func Parse(r io.Reader) (*Scraper, error) {
	return nil, nil
}

func FromPath(path string) (*Scraper, error) {
	return &Scraper{}, nil
}
