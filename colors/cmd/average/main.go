package main

import (
	"fmt"
	"log"

	"github.com/grokify/mogo/image/colors"
	"github.com/grokify/mogo/image/imageutil"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	Filename string `short:"f" long:"filename" description:"Filename" required:"true"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := imageutil.ReadImageFile(opts.Filename)
	if err != nil {
		log.Fatal(err)
	}
	avgClr := colors.ColorAverageImage(img)
	fmt.Printf("COLOR [%s]\n", colors.ColorToHex(avgClr))

	if 1 == 0 {
		img2 := imageutil.AddBorderAverageColor(img, 100)
		err = imageutil.WriteFileJPEG("_with_border.jpg", img2, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("DONE")
}

/*
func AddBorderAverageColor(img image.Image, width uint) (image.Image, error) {
	imgRGBA := imageutil.ImageToRGBA(img)
	if width < 1 {
		return imgRGBA, errors.New("zero width border")
	}
	avgClr := colors.ColorAverageImage(imgRGBA)
	return imageutil.AddBorder(imageutil.ImageToRGBA(img), avgClr, width), nil
}
*/
