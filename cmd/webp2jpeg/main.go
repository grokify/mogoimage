package main

import (
	"fmt"
	"image/jpeg"
	"log"

	"github.com/grokify/goimage"
	"github.com/grokify/mogo/image/imageutil"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	Input       string `short:"i" long:"input" description:"Webp URL or filepath" required:"true"`
	Output      string `short:"o" long:"output" description:"JPEG filepath" required:"true"`
	JPEGQuality uint   `short:"q" long:"quality" description:"JPEG quality" default:"80"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	img, ext, err := goimage.ReadImageAny(opts.Input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("GOT TYPE [%s]\n", ext)

	if opts.JPEGQuality == 0 {
		opts.JPEGQuality = jpeg.DefaultQuality
	}

	i2 := imageutil.Image{Image: img}

	err = i2.WriteJPEGFile(opts.Output,
		&imageutil.JPEGEncodeOptions{
			Options: &jpeg.Options{Quality: int(opts.JPEGQuality)}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
}
