package mogoimage

import (
	"fmt"
	"image"
	"image/color"

	"github.com/grokify/mogo/image/colors"

	micolors "github.com/grokify/mogoimage/colors"
)

// CropImage takes an image and crops it to the specified rectangle. `CropImage`
// is from: https://stackoverflow.com/a/63256403.
func CropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	// img is an Image interface. This checks if the underlying value has a
	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}

func CropImageColor(img image.Image, tolerance float64, remove ...color.Color) (image.Image, error) {
	cols := ColumnsFilter(img, tolerance, remove...)
	if len(cols) == 0 {
		return img, nil
	}

	newRect := img.Bounds().Bounds()

	if int(cols[0]) == img.Bounds().Min.X {
		trimLeft := 0
		for x := 1; x < len(cols); x++ {
			if cols[x] == cols[x-1]+1 {
				trimLeft = int(cols[x])
			} else {
				trimLeft = int(cols[x-1]) + 1
				break
			}
		}
		newRect.Min.X = trimLeft
	}
	if int(cols[len(cols)-1]) == img.Bounds().Max.X-1 {
		trimRight := 0
		for x := len(cols) - 1; x >= 0; x-- {
			if x == 0 {
				trimRight = int(cols[0])
			} else if cols[x] == cols[x-1]+1 {
				trimRight = int(cols[x-1])
			} else {
				break
			}
		}
		newRect.Max.X = trimRight
	}
	if newRect.Eq(img.Bounds()) {
		return img, nil
	}
	return CropImage(img, newRect)
}

// ColumnsFilter returns a list of column indexes that matches the wanted colors
// within the provided tolerane.
func ColumnsFilter(img image.Image, tolerance float64, want ...color.Color) []uint {
	cols := []uint{}
	if len(want) == 0 {
		return cols
	}
	wantColorsUnique := micolors.ColorsDistance(colors.SliceUnique(want))
	minPt := img.Bounds().Min
	maxPt := img.Bounds().Max
	for x := minPt.X; x < maxPt.X; x++ {
		colColors := []color.Color{}
		for y := minPt.Y; y < maxPt.Y; y++ {
			colColors = append(colColors, img.At(x, y))
		}
		colColorsUniqueDist := micolors.ColorsDistance(colors.SliceUnique(colColors))
		if wantColorsUnique.MatchBest(tolerance, colColorsUniqueDist...) {
			cols = append(cols, uint(x))
		}
	}
	return cols
}
