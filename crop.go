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
func CropImage(img image.Image, retain image.Rectangle) (image.Image, error) {
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

	return simg.SubImage(retain), nil
}

func CropImageColor(img image.Image, tolerance float64, remove ...color.Color) (image.Image, error) {
	cols := ColumnsFilter(img, tolerance, remove...)
	if len(cols) == 0 {
		return img, nil
	}

	newRect := img.Bounds().Bounds().Canon()

	if int(cols[0]) == img.Bounds().Canon().Min.X {
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
	if int(cols[len(cols)-1]) == img.Bounds().Canon().Max.X-1 {
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
	if newRect.Eq(img.Bounds().Canon()) {
		return img, nil
	}
	return CropImage(img, newRect)
}

func CropImageColorCaption(img image.Image, paddingPct, tolerance float64, remove ...color.Color) (image.Image, error) {
	rows := RowsFilterCaption(img, paddingPct, tolerance, remove...)
	if len(rows) == 0 {
		return img, nil
	}

	newRect := img.Bounds().Canon().Bounds()

	if int(rows[len(rows)-1]) == img.Bounds().Canon().Max.Y-1 {
		trimBottom := 0
		for y := len(rows) - 1; y >= 0; y-- {
			if y == 0 {
				trimBottom = int(rows[0])
			} else if rows[y] == rows[y-1]+1 {
				trimBottom = int(rows[y-1])
			} else {
				break
			}
		}
		validatedTrimBottom := 0
		for y := trimBottom; y < img.Bounds().Canon().Max.Y; y++ {
			if RowMatch(img, y, tolerance, remove...) {
				validatedTrimBottom = y
				break
			}
		}
		newRect.Max.Y = validatedTrimBottom
	}
	if newRect.Eq(img.Bounds().Canon()) {
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
	minPt := img.Bounds().Canon().Min
	maxPt := img.Bounds().Canon().Max
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

// RowsFilter returns a list of row indexes that matches the wanted colors
// within the provided tolerane.
func RowsFilter(img image.Image, tolerance float64, want ...color.Color) []uint {
	rows := []uint{}
	if len(want) == 0 {
		return rows
	}
	wantColorsUnique := micolors.ColorsDistance(colors.SliceUnique(want))
	minPt := img.Bounds().Canon().Min
	maxPt := img.Bounds().Canon().Max
	for y := minPt.Y; y < maxPt.Y; y++ {
		rowColors := []color.Color{}
		for x := minPt.X; x < maxPt.X; x++ {
			rowColors = append(rowColors, img.At(x, y))
		}
		colColorsUniqueDist := micolors.ColorsDistance(colors.SliceUnique(rowColors))
		if wantColorsUnique.MatchBest(tolerance, colColorsUniqueDist...) {
			rows = append(rows, uint(y))
		}
	}
	return rows
}

// RowMatch checks to see if an image row matches the want colors.
func RowMatch(img image.Image, rowIdx int, tolerance float64, want ...color.Color) bool {
	if rowIdx <= 0 {
		return false
	}
	wantColorsUnique := micolors.ColorsDistance(colors.SliceUnique(want))
	minPt := img.Bounds().Canon().Min
	maxPt := img.Bounds().Canon().Max
	if rowIdx >= maxPt.Y {
		return false
	}
	y := rowIdx
	rowColors := []color.Color{}
	for x := minPt.X; x < maxPt.X; x++ {
		rowColors = append(rowColors, img.At(x, y))
	}
	colColorsUniqueDist := micolors.ColorsDistance(colors.SliceUnique(rowColors))
	return wantColorsUnique.MatchBest(tolerance, colColorsUniqueDist...)
}

// RowsFilterCaption returns a list of column indexes that matches the wanted colors
// within the provided tolerane.
func RowsFilterCaption(img image.Image, paddingPct, tolerance float64, want ...color.Color) []uint {
	rows := []uint{}
	if len(want) == 0 {
		return rows
	}
	wantColorsUnique := micolors.ColorsDistance(colors.SliceUnique(want))
	minPt := img.Bounds().Canon().Min
	maxPt := img.Bounds().Canon().Max
	for y := minPt.Y; y < maxPt.Y; y++ {
		rowColors := []color.Color{}
		for x := minPt.X; x < maxPt.X; x++ {
			rowColors = append(rowColors, img.At(x, y))
		}
		if RowMatchCaption(paddingPct, tolerance, wantColorsUnique, rowColors) {
			rows = append(rows, uint(y))
		}
	}
	return rows
}

// RowCaptionMatch matches a row with a caption. paddingPct is what is considered for matching.
// padding is left or right. Padding 0.5 or greater means matching the entire row. Padding
// <= 0 means no matching is necessary to succeed. Negative padding is not supported and converted
// to 0.
func RowMatchCaption(paddingPct, tolerance float64, wantColorsUnique micolors.ColorsDistance, candidate []color.Color) bool {
	if len(wantColorsUnique) == 0 || len(candidate) == 0 {
		return false
	} else if paddingPct <= 0 {
		return true
	}
	if paddingPct >= 0.5 {
		for _, cx := range candidate {
			if !wantColorsUnique.MatchBest(tolerance, cx) {
				return false
			}
		}
		return true
	}
	candidateCount := len(candidate)
	paddingCount := int(float64(candidateCount) * paddingPct)
	// no padding to match
	if paddingCount == 0 {
		return true
	}
	// try left
	for x := 0; x <= paddingCount; x++ {
		cx := candidate[x]
		if !wantColorsUnique.MatchBest(tolerance, cx) {
			return false
		}
	}
	// try right
	// len = 5, p = 1  x = 4, succ, test >= 5-1 = 4; succ
	// len = 5, p = 1, x = 3, fail, test >= 5-1 = 4; fail
	for x := len(candidate) - 1; x >= len(candidate)-paddingCount; x-- {
		cx := candidate[x]
		if !wantColorsUnique.MatchBest(tolerance, cx) {
			return false
		}
	}
	return true
}
