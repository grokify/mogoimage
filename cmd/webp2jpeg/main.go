package main

import (
	"fmt"
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

	err = imageutil.WriteFileJPEG(opts.Output, img, int(opts.JPEGQuality))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
}
