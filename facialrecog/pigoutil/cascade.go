package pigoutil

import (
	_ "embed"
)

//go:embed "cascade/facefinder"
var cascadeFaceFinderEmbed []byte

//go:embed "cascade/puploc"
var cascadePupilFinderEmbed []byte

func CascadeFaceFinder() []byte {
	return cascadeFaceFinderEmbed
}

func CascadePupilFinder() []byte {
	return cascadePupilFinderEmbed
}
