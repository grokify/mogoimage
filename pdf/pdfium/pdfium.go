package pdfium

import (
	"errors"
	"image"
	"os"

	render "github.com/brunsgaard/go-pdfium-render"
	"github.com/grokify/mogo/image/imageutil"
)

type Metadata struct {
	Width  uint
	Height uint
}

func ImageMetadata(img image.Image) Metadata {
	m := Metadata{
		Width:  uint(img.Bounds().Dx()),
		Height: uint(img.Bounds().Dy()),
	}
	return m
}

func ConvertPDFPageToPNG(srcFilename, dstFilename string, pgIndex, dpi, minWidth uint) error {
	img, err := ReadPDFPageImage(srcFilename, pgIndex, dpi)
	if err != nil {
		return err
	}
	// m := ImageMetadata(img)
	// fmtutil.PrintJSON(m)
	i2 := imageutil.ResizeMin(minWidth, 0, img, imageutil.ScalerBest())
	i3 := imageutil.Image{Image: i2}
	return i3.WritePNGFile(dstFilename)
}

func ReadPDFPageImage(filename string, pgIndex, dpi uint) (*image.RGBA, error) {
	if dpi == 0 {
		dpi = 72
	}
	d, err := ReadPDF(filename)
	if err != nil {
		return nil, err
	}
	if int(pgIndex) >= d.GetPageCount() {
		return nil, errors.New("page index does not exist")
	}
	// fmt.Printf("Page_COUNT (%d)(%d)\n", d.GetPageCount(), pgIndex)
	// fmt.Printf("DPI (%v)\n", dpi)
	return d.RenderPage(int(pgIndex), int(dpi)), nil
}

func ReadPDF(filename string) (*render.Document, error) {
	if b, err := os.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return render.NewDocument(&b)
	}
}
