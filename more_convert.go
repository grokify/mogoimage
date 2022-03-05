package goimage

import (
	"image"

	"github.com/andybons/gogif"
)

func ToPalettedMedianCut(src image.Image) *image.Paletted {
	if v, ok := src.(*image.Paletted); ok {
		return v
	}
	pimg := image.NewPaletted(src.Bounds(), nil)
	quantizer := gogif.MedianCutQuantizer{NumColor: 256}
	quantizer.Quantize(pimg, src.Bounds(), src, image.ZP)
	return pimg
}
