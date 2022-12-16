package tui

type state int

const (
	stateExtensionSelect state = iota + 1
	stateExtensionConfig
	stateSearch
	stateSearchResults
	stateLayer
	statePrepare
	stateStreamOrDownloadSelection
	stateStream
	stateDownload
)