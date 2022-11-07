package colors

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/grokify/mogo/image/colors"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	DistanceDefault = "cie76"
	DistanceGood    = "cie76"
	DistanceBetter  = "cie94"
	DistanceBest    = "cie2k"

	ToleranceNone        = 0.0
	ToleranceBestDefault = 0.05
)

func Distance(alg string, c1, c2 color.Color) (float64, error) {
	switch strings.ToLower(strings.TrimSpace(alg)) {
	case DistanceGood:
		return DistanceCIE76(c1, c2), nil
	case DistanceBetter:
		return DistanceCIE94(c1, c2), nil
	case DistanceBest:
		return DistanceCIE2K(c1, c2), nil
	}
	return 0.0, fmt.Errorf("algorithm not known [%s]", alg)
}

func MustDistance(alg string, c1, c2 color.Color) float64 {
	d, err := Distance(alg, c1, c2)
	if err != nil {
		panic(fmt.Errorf("algorithm not known [%s]", alg))
	}
	return d
}

func Distances(alg string, comp color.Color, c []color.Color) ([]float64, error) {
	dists := []float64{}
	for _, cx := range c {
		dist, err := Distance(alg, comp, cx)
		if err != nil {
			return dists, err
		}
		dists = append(dists, dist)
	}
	return dists, nil
}

type ColorDistance struct {
	Color    color.Color
	ColorHex string
	Distance float64
}

func DistancesMore(alg string, comp color.Color, c []color.Color) ([]ColorDistance, error) {
	dists := []ColorDistance{}
	for _, cx := range c {
		dist, err := Distance(alg, comp, cx)
		if err != nil {
			return dists, err
		}
		dists = append(dists, ColorDistance{
			Color:    cx,
			ColorHex: colors.ColorToHex(cx),
			Distance: dist})
	}
	return dists, nil
}

func DistancesMatrix(alg string, comp color.Color, c [][]color.Color) ([][]float64, error) {
	distsM := [][]float64{}
	for _, cx := range c {
		dists, err := Distances(alg, comp, cx)
		if err != nil {
			return distsM, err
		}
		distsM = append(distsM, dists)
	}
	return distsM, nil
}

func DistanceCIE2K(c1, c2 color.Color) float64 {
	return ColorfulColor(c1).DistanceCIEDE2000(ColorfulColor(c2))
}

func DistanceCIE94(c1, c2 color.Color) float64 {
	return ColorfulColor(c1).DistanceCIE94(ColorfulColor(c2))
}

func DistanceCIE76(c1, c2 color.Color) float64 {
	return ColorfulColor(c1).DistanceCIE76(ColorfulColor(c2))
}

func ColorfulColor(c color.Color) colorful.Color {
	r, g, b, _ := c.RGBA()
	return colorful.Color{
		R: float64(r/256) / 255.0,
		G: float64(g/256) / 255.0,
		B: float64(b/256) / 255.0}
	/*
		return colorful.Color{
			R: float64(clr.R) / 255.0,
			G: float64(clr.G) / 255.0,
			B: float64(clr.B) / 255.0}c2
	*/
}

type ColorsDistance []color.Color

func (cd ColorsDistance) Unique() ColorsDistance {
	return colors.SliceUnique(cd)
}

func (cd ColorsDistance) MatchBest(tolerance float64, c ...color.Color) bool {
	if len(c) == 0 {
		return false
	}
	absTolerance := AbsFloat64(tolerance)
	for _, ctest := range c {
		for _, cx := range cd {
			if AbsFloat64(DistanceCIE2K(ctest, cx)) > absTolerance {
				return false
			}
		}
	}
	return true
}

func AbsFloat64(v float64) float64 {
	if v < 0 {
		return -1.0 * v
	}
	return v
}
