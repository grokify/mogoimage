package mupdf

import (
	"fmt"
	"image"

	fitz "github.com/gen2brain/go-fitz"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/image/imageutil"
	"golang.org/x/image/draw"
)

/*

use `brew intall mupdf-tools` for `go-fitz` muPDF dependencies.

*/

func ConvertFilePageToImage(filename string, pageIndex uint) (image.Image, error) {
	doc, err := fitz.New(filename)
	if err != nil {
		return nil, errorsutil.Wrapf(err, "error opening path with `go-fitz` (%s)", filename)
	} else {
		defer doc.Close()
	}
	if int(pageIndex) >= doc.NumPage() {
		return nil, fmt.Errorf("requested page index out of range: requested (%d) page count (%d)", pageIndex, doc.NumPage())
	} else if img, err := doc.Image(int(pageIndex)); err != nil {
		return img, errorsutil.Wrapf(err, "error using go-fitz.Document.Image(\"%d\")", pageIndex)
	} else {
		return img, nil
	}
}

func ConvertFilePageToPNGFile(srcPath, outPath string, pageIndex, width, height uint, scaler draw.Scaler) error {
	if img, err := ConvertFilePageToImage(srcPath, pageIndex); err != nil {
		return err
	} else {
		if width > 0 || height > 0 {
			img = imageutil.Resize(width, height, img, scaler)
		}
		i2 := imageutil.Image{Image: img}
		return i2.WritePNGFile(outPath)
	}
}
