package app

import "github.com/awirix/awirix/version"

const (
	Name   = "awirix"
	Prefix = Name + "_"
)

var (
	Version = version.MustParse("0.0.1")
)
