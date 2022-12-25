package tui

//go:generate enumer -type=variant -trimprefix=variant -transform=kebab-case
type variant int

const (
	variantStream variant = iota + 1
	variantDownload
)
