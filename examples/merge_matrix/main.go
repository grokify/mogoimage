package main

import (
	"fmt"
	"log"

	imageutil "github.com/grokify/mogoimage"
)

func main() {
	urlsmatrix := [][]string{
		{"https://raw.githubusercontent.com/grokify/mogoimage/master/read_testdata/gopher_appengine_color.jpg"},
		{
			"https://raw.githubusercontent.com/grokify/mogoimage/master/read_testdata/gopher_color.jpg",
			"https://raw.githubusercontent.com/grokify/mogoimage/master/read_testdata/gopher_color.jpg"},
	}

	outfile := "_merged.jpg"

	matrix, err := imageutil.MatrixRead(urlsmatrix)
	if err != nil {
		log.Fatal(err)
	}

	err = imageutil.WriteFileJPEG(
		outfile, matrix.Merge(true, true), imageutil.JPEGQualityMax)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%v]\n", outfile)

	fmt.Println("DONE")
}
