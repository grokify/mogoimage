package main

import (
	"fmt"
	"image"
	"log"
	"strings"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/image/imageutil"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	Images     []string `short:"i" description:"A slice of image files" required:"true"`
	Outputfile string   `short:"o" description:"Output file" required:"true"`
	Quality    int      `short:"q" description:"Quality"`
	Verbose    []bool   `short:"v" description:"Verbose logging"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.MustPrintJSON(opts.Images)

	imSet, err := imageutil.NewImageSetFiles(opts.Images)
	if err != nil {
		log.Fatal(err)
	}

	if len(opts.Verbose) > 0 {
		fmtutil.MustPrintJSON(imSet)
		fmtutil.MustPrintJSON(imSet.Stats())
	}

	imgs := []image.Image{}
	for _, im := range imSet.ImageMetas {
		imgs = append(imgs, im.Image)
	}

	//merged := imageutil.MergeHorizontalRGBA(imSet)
	merged := imageutil.MergeYSameX(imgs, true)

	filename := strings.TrimSpace(opts.Outputfile)
	if len(filename) == 0 {
		filename = "merged.jpg"
	}

	err = imageutil.WriteFileJPEG(filename, merged, opts.Quality)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("MERGED: %v\n", strings.Join(opts.Images, ", "))
	fmt.Printf("WROTE %s\n", filename)

	fmt.Println("DONE")
}
