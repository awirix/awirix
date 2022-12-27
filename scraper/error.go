package scraper

import "errors"

var (
	ErrMissingHandler = errors.New("missing handler")
	ErrMissingTitle   = errors.New("missing title")
)
