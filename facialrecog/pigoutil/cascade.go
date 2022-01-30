package pigoutil

import (
	_ "embed"
	"fmt"
	"strings"
)

const (
	PathCascadeFacefinder = "cascade/facefinder"
	PathCascadePupil      = "cascade/puploc"
)

//go:embed "cascade/facefinder"
var cascadeFaceFinderEmbed []byte

//go:embed "cascade/puploc"
var cascadePupilFinderEmbed []byte

func ReadCascadeFaceFinder() []byte {
	return cascadeFaceFinderEmbed
}

func ReadCascadePupilFinder() []byte {
	return cascadePupilFinderEmbed
}

func ReadCascade(cascadeType string) ([]byte, error) {
	cascadeType = strings.ToLower(strings.TrimSpace(cascadeType))
	if strings.Contains(cascadeType, "facefinder") ||
		strings.Contains(cascadeType, "face") {
		return cascadeFaceFinderEmbed, nil
	} else if strings.Contains(cascadeType, "puploc") ||
		strings.Contains(cascadeType, "pupil") {
		return cascadePupilFinderEmbed, nil
	}
	return []byte{}, fmt.Errorf("cascade not found [%s]", cascadeType)
}
