package colors

import (
	"image/color"
	"testing"

	mogocolors "github.com/grokify/mogo/image/colors"
)

var colorDistanceTests = []struct {
	c1         color.Color
	c2         color.Color
	distGood   float64
	distBetter float64
	distBest   float64
}{
	{color.Black, color.Black, 0.0, 0.0, 0.0},
	{color.Black, color.White, 1.0000000100258313, 1.0000000100258313, 1.0000000103996842},
	{color.White, color.White, 0.0, 0.0, 0.0},
	{color.White, mogocolors.Red, 1.145380009323503, 1.1447724245372048, 0.45815560529554045},
	{color.White, mogocolors.Green, 1.2041092506970081, 1.203350603752062, 0.3327081037940642},
	{color.White, mogocolors.Blue, 1.4995965106584739, 1.498836259416328, 0.639087533748455},
}

func TestColorDistances(t *testing.T) {
	for _, tt := range colorDistanceTests {
		distGood, err := Distance(DistanceGood, tt.c1, tt.c2)
		if err != nil {
			t.Errorf("Distance(\"%s\", ...): error [%v]", DistanceGood, err.Error())
		}
		if distGood != tt.distGood {
			t.Errorf("Distance(\"%s\", ...): want [%v] got [%v]", DistanceGood, tt.distGood, distGood)
		}
		distBetter, err := Distance(DistanceBetter, tt.c1, tt.c2)
		if err != nil {
			t.Errorf("Distance\"%s\", (...): error [%v]", DistanceBetter, err.Error())
		}
		if distBetter != tt.distBetter {
			t.Errorf("Distance(\"%s\", ...): want [%v] got [%v]", DistanceBetter, tt.distBetter, distBetter)
		}
		distBest, err := Distance(DistanceBest, tt.c1, tt.c2)
		if err != nil {
			t.Errorf("Distance\"%s\", (...): error [%v]", DistanceBest, err.Error())
		}
		if distBest != tt.distBest {
			t.Errorf("Distance(\"%s\", ...): want [%v] got [%v]", DistanceBest, tt.distBest, distBest)
		}
	}
}
