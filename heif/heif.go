package heif

import (
	"errors"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/adrium/goheif"
	goheifheif "github.com/adrium/goheif/heif"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/image/imageutil"
	"github.com/grokify/mogo/io/ioutil"
	"github.com/grokify/mogo/os/osutil"
	"github.com/grokify/mogo/path/filepathutil"
)

var rxHEICExtension = regexp.MustCompile(`(?i)\.heic$`)

// WriteJPEGDir uses `imageutil.JPEGEncodeOptions` properties `WriteExtension`, and `Options`.
func WriteJPEGDir(outdir, srcdir string, rxHEICFilename *regexp.Regexp, force, recurse bool, opt *imageutil.JPEGEncodeOptions) error {
	if opt == nil {
		opt = &imageutil.JPEGEncodeOptions{}
	}
	if rxHEICFilename == nil {
		rxHEICFilename = rxHEICExtension
	}
	entries, err := osutil.ReadDirMore(srcdir, rxHEICFilename, false, true, false)
	if err != nil {
		return err
	}
	ext := opt.WriteExtensionOrDefault()

	for _, e := range entries {
		heicfilename := filepath.Join(srcdir, e.Name())
		jpegfilename := filepath.Join(outdir, filepathutil.TrimExt(e.Name())+ext)
		if !force {
			if ok, err := osutil.IsFile(jpegfilename, true); err == nil && ok {
				continue
			}
		}
		if err := WriteJPEGFile(jpegfilename, heicfilename, opt.Options); err != nil {
			return err
		}
	}

	if recurse {
		sdirs, err := osutil.ReadDirMore(srcdir, nil, true, false, false) // call again to avoid using regexp match on subdirs.
		if err != nil {
			return err
		}
		for _, sdir := range sdirs {
			if err := WriteJPEGDir(filepath.Join(outdir, sdir.Name()), filepath.Join(srcdir, sdir.Name()), rxHEICFilename, force, recurse, opt); err != nil {
				return err
			}
		}
	}
	return nil
}

func WriteJPEGFile(jpegName, heicName string, opt *jpeg.Options) error {
	fi, err := os.Open(heicName)
	if err != nil {
		return err
	}
	defer fi.Close()

	fo, err := os.OpenFile(jpegName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer fo.Close()
	return WriteJPEG(fo, fi, opt)
}

func WriteJPEG(w io.Writer, r ioutil.AtReader, opt *jpeg.Options) error {
	exif, err := goheif.ExtractExif(r)
	if err != nil {
		if errors.Is(err, goheifheif.ErrNoEXIF) {
			exif = []byte{}
		} else {
			return errorsutil.Wrap(err, "goheif.ExtractExif() error: ")
		}
	}

	img, err := goheif.Decode(r)
	if err != nil {
		return errorsutil.Wrap(err, "failed to parse HEIC io.Reader using goheif.Decode()")
	}

	err = imageutil.Image{Image: img}.WriteJPEG(w, &imageutil.JPEGEncodeOptions{
		Options: opt,
		Exif:    exif,
	})
	if err != nil {
		return errorsutil.Wrap(err, "failed to write JPEG to io.Writer using imageutil.WriteJPEG()")
	}
	return nil
}
