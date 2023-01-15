package pigoutil

import (
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
	"github.com/esimov/pigo/utils"
	"github.com/fogleman/gg"
	"github.com/grokify/mogo/net/http/httputilmore"
	"golang.org/x/term"
)

const (
	// pipeName is the file name that indicates stdin/stdout is being used.
	pipeName = "-"
	// markerRectangle - use rectangle as face detection marker
	markerRectangle string = "rect"
	// markerCircle - use circle as face detection marker
	markerCircle string = "circle"
	// markerEllipse - use ellipse as face detection marker
	markerEllipse string = "ellipse"

	// message colors
	successColor = "\x1b[92m"
	errorColor   = "\x1b[31m"
	defaultColor = "\x1b[0m"

	perturb = 63

	ExtensionPNG  = ".png"
	ExtensionJPG  = ".jpg"
	ExtensionJPEG = ".jpeg"
)

func eyeCascades() []string {
	return []string{"lp46", "lp44", "lp42", "lp38", "lp312"}
}

func mouthCascades() []string {
	return []string{"lp93", "lp84", "lp82", "lp81"}
}

func FileExtensionSupported() []string {
	return []string{ExtensionJPG, ExtensionJPEG, ExtensionPNG}
}

type FaceDetector struct {
	Angle        float64
	CascadeFile  string
	Destination  string
	MinSize      int
	MaxSize      int
	ShiftFactor  float64
	ScaleFactor  float64
	IouThreshold float64
	Puploc       string
	Flploc       string
	MarkDetEyes  bool
	dc           *gg.Context
	det          *FaceDetector
	plc          *pigo.PuplocCascade
	flpcs        map[string][]*pigo.FlpCascade
	imgParams    *pigo.ImageParams
}

type FaceDetectorVars struct {
	//dc        *gg.Context
	det   *FaceDetector
	plc   *pigo.PuplocCascade
	flpcs map[string][]*pigo.FlpCascade
	//imgParams *pigo.ImageParams
}

// DetectFaces run the detection algorithm over the provided source image.
func (fd *FaceDetector) DetectFaces(source string) ([]pigo.Detection, error) {
	var srcFile io.Reader

	// Check if source path is a local image or URL.
	if utils.IsValidUrl(source) {
		src, err := utils.DownloadImage(source)
		if err != nil {
			return nil, err
		}
		// Close and remove the generated temporary file.
		defer src.Close()
		defer os.Remove(src.Name())

		img, err := os.Open(src.Name())
		if err != nil {
			return nil, err
		}
		srcFile = img
	} else {
		if source == pipeName {
			if term.IsTerminal(int(os.Stdin.Fd())) {
				log.Fatalln("`-` should be used with a pipe for stdin")
			}
			srcFile = os.Stdin
		} else {
			file, err := os.Open(source)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			srcFile = file
		}
	}

	src, err := pigo.DecodeImage(srcFile)
	if err != nil {
		return nil, err
	}

	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	fd.dc = gg.NewContext(cols, rows)
	fd.dc.DrawImage(src, 0, 0)

	fd.imgParams = &pigo.ImageParams{
		Pixels: pixels,
		Rows:   rows,
		Cols:   cols,
		Dim:    cols,
	}

	cParams := pigo.CascadeParams{
		MinSize:     fd.MinSize,
		MaxSize:     fd.MaxSize,
		ShiftFactor: fd.ShiftFactor,
		ScaleFactor: fd.ScaleFactor,
		ImageParams: *fd.imgParams,
	}

	// cascadeFile, err := ioutil.ReadFile(fd.CascadeFile)
	cascadeFileBytes, err := ReadCascade(fd.CascadeFile)
	if err != nil {
		return nil, err
	}

	// ContentType New
	contentType := http.DetectContentType(cascadeFileBytes)
	if contentType != httputilmore.ContentTypeAppOctetStream {
		return nil, errors.New("the provided cascade classifier is not valid")
	}

	p := pigo.NewPigo()

	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err := p.Unpack(cascadeFileBytes)
	if err != nil {
		return nil, err
	}

	if len(fd.Puploc) > 0 {
		pl := pigo.NewPuplocCascade()
		//cascade, err := ioutil.ReadFile(fd.Puploc)
		cascade, err := ReadCascade(fd.Puploc)
		if err != nil {
			return nil, err
		}
		plc, err := pl.UnpackCascade(cascade)
		if err != nil {
			return nil, err
		}
		fd.plc = plc

		if len(fd.Flploc) > 0 {
			flpcs, err := pl.ReadCascadeDir(fd.Flploc)
			if err != nil {
				return nil, err
			}
			fd.flpcs = flpcs
		}
	}

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	faces := classifier.RunCascade(cParams, fd.Angle)

	// Calculate the intersection over union (IoU) of two clusters.
	faces = classifier.ClusterDetections(faces, fd.IouThreshold)

	return faces, nil
}

// DrawFaces marks the detected faces with the marker type defined as parameter (rectangle|circle|ellipse).
func (fd *FaceDetector) DrawFaces(faces []pigo.Detection, marker string) ([]detection, error) {
	var qThresh float32 = 5.0

	var (
		detections     = make([]detection, 0, len(faces))
		eyesCoords     = make([]coord, 0, len(faces))
		landmarkCoords = make([]coord, 0, len(faces))
		puploc         *pigo.Puploc
	)

	for _, face := range faces {
		if face.Q > qThresh {
			switch marker {
			case markerRectangle:
				fd.dc.DrawRectangle(float64(face.Col-face.Scale/2),
					float64(face.Row-face.Scale/2),
					float64(face.Scale),
					float64(face.Scale),
				)
			case markerCircle:
				fd.dc.DrawArc(
					float64(face.Col),
					float64(face.Row),
					float64(face.Scale/2),
					0,
					2*math.Pi,
				)
			case markerEllipse:
				fd.dc.DrawEllipse(
					float64(face.Col),
					float64(face.Row),
					float64(face.Scale)/2,
					float64(face.Scale)/1.6,
				)
			}
			faceCoord := &coord{
				Col:   face.Row - face.Scale/2,
				Row:   face.Col - face.Scale/2,
				Scale: face.Scale,
			}

			fd.dc.SetLineWidth(2.0)
			fd.dc.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{R: 255, G: 0, B: 0, A: 255}))
			fd.dc.Stroke()

			if len(fd.Puploc) > 0 && face.Scale > 50 {
				rect := image.Rect(
					face.Col-face.Scale/2,
					face.Row-face.Scale/2,
					face.Col+face.Scale/2,
					face.Row+face.Scale/2,
				)
				rows, cols := rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y
				ctx := gg.NewContext(rows, cols)
				faceZone := ctx.Image()

				// left eye
				puploc = &pigo.Puploc{
					Row:      face.Row - int(0.075*float32(face.Scale)),
					Col:      face.Col - int(0.175*float32(face.Scale)),
					Scale:    float32(face.Scale) * 0.25,
					Perturbs: perturb,
				}
				leftEye := fd.plc.RunDetector(*puploc, *fd.imgParams, fd.Angle, false)
				if leftEye.Row > 0 && leftEye.Col > 0 {
					if fd.Angle > 0 {
						drawEyeDetectionMarker(ctx,
							float64(cols/2-(face.Col-leftEye.Col)),
							float64(rows/2-(face.Row-leftEye.Row)),
							float64(leftEye.Scale),
							color.RGBA{R: 255, G: 0, B: 0, A: 255},
							fd.MarkDetEyes,
						)
						angle := (fd.Angle * 180) / math.Pi
						rotated := imaging.Rotate(faceZone, 2*angle, color.Transparent)
						final := imaging.FlipH(rotated)

						fd.dc.DrawImage(final, face.Col-face.Scale/2, face.Row-face.Scale/2)
					} else {
						drawEyeDetectionMarker(fd.dc,
							float64(leftEye.Col),
							float64(leftEye.Row),
							float64(leftEye.Scale),
							color.RGBA{R: 255, G: 0, B: 0, A: 255},
							fd.MarkDetEyes,
						)
					}
					eyesCoords = append(eyesCoords, coord{
						Col:   leftEye.Row,
						Row:   leftEye.Col,
						Scale: int(leftEye.Scale),
					})
				}

				// right eye
				puploc = &pigo.Puploc{
					Row:      face.Row - int(0.075*float32(face.Scale)),
					Col:      face.Col + int(0.185*float32(face.Scale)),
					Scale:    float32(face.Scale) * 0.25,
					Perturbs: perturb,
				}

				rightEye := fd.plc.RunDetector(*puploc, *fd.imgParams, fd.Angle, false)
				if rightEye.Row > 0 && rightEye.Col > 0 {
					if fd.Angle > 0 {
						drawEyeDetectionMarker(ctx,
							float64(cols/2-(face.Col-rightEye.Col)),
							float64(rows/2-(face.Row-rightEye.Row)),
							float64(rightEye.Scale),
							color.RGBA{R: 255, G: 0, B: 0, A: 255},
							fd.MarkDetEyes,
						)
						// convert radians to angle
						angle := (fd.Angle * 180) / math.Pi
						rotated := imaging.Rotate(faceZone, 2*angle, color.Transparent)
						final := imaging.FlipH(rotated)

						fd.dc.DrawImage(final, face.Col-face.Scale/2, face.Row-face.Scale/2)
					} else {
						drawEyeDetectionMarker(fd.dc,
							float64(rightEye.Col),
							float64(rightEye.Row),
							float64(rightEye.Scale),
							color.RGBA{R: 255, G: 0, B: 0, A: 255},
							fd.MarkDetEyes,
						)
					}
					eyesCoords = append(eyesCoords, coord{
						Col:   rightEye.Row,
						Row:   rightEye.Col,
						Scale: int(rightEye.Scale),
					})
				}

				if len(fd.Flploc) > 0 {
					for _, eye := range eyeCascades() {
						for _, flpc := range fd.flpcs[eye] {
							flp := flpc.GetLandmarkPoint(leftEye, rightEye, *fd.imgParams, perturb, false)
							if flp.Row > 0 && flp.Col > 0 {
								drawEyeDetectionMarker(fd.dc,
									float64(flp.Col),
									float64(flp.Row),
									float64(flp.Scale*0.5),
									color.RGBA{R: 0, G: 0, B: 255, A: 255},
									false,
								)
								landmarkCoords = append(landmarkCoords, coord{
									Col:   flp.Row,
									Row:   flp.Col,
									Scale: int(flp.Scale),
								})
							}

							flp = flpc.GetLandmarkPoint(leftEye, rightEye, *fd.imgParams, perturb, true)
							if flp.Row > 0 && flp.Col > 0 {
								drawEyeDetectionMarker(fd.dc,
									float64(flp.Col),
									float64(flp.Row),
									float64(flp.Scale*0.5),
									color.RGBA{R: 0, G: 0, B: 255, A: 255},
									false,
								)
								landmarkCoords = append(landmarkCoords, coord{
									Col:   flp.Row,
									Row:   flp.Col,
									Scale: int(flp.Scale),
								})
							}
						}
					}

					for _, mouth := range mouthCascades() {
						for _, flpc := range fd.flpcs[mouth] {
							flp := flpc.GetLandmarkPoint(leftEye, rightEye, *fd.imgParams, perturb, false)
							if flp.Row > 0 && flp.Col > 0 {
								drawEyeDetectionMarker(fd.dc,
									float64(flp.Col),
									float64(flp.Row),
									float64(flp.Scale*0.5),
									color.RGBA{R: 0, G: 0, B: 255, A: 255},
									false,
								)
								landmarkCoords = append(landmarkCoords, coord{
									Col:   flp.Row,
									Row:   flp.Col,
									Scale: int(flp.Scale),
								})
							}
						}
					}
					flp := fd.flpcs["lp84"][0].GetLandmarkPoint(leftEye, rightEye, *fd.imgParams, perturb, true)
					if flp.Row > 0 && flp.Col > 0 {
						drawEyeDetectionMarker(fd.dc,
							float64(flp.Col),
							float64(flp.Row),
							float64(flp.Scale*0.5),
							color.RGBA{R: 0, G: 0, B: 255, A: 255},
							false,
						)
						landmarkCoords = append(landmarkCoords, coord{
							Col:   flp.Row,
							Row:   flp.Col,
							Scale: int(flp.Scale),
						})
					}
				}
			}
			detections = append(detections, detection{
				FacePoints:     *faceCoord,
				EyePoints:      eyesCoords,
				LandmarkPoints: landmarkCoords,
			})
		}
	}
	return detections, nil
}

func (fd *FaceDetector) EncodeImage(dst io.Writer) error {
	var err error
	img := fd.dc.Image()

	switch dst.(type) {
	case *os.File:
		ext := filepath.Ext(dst.(*os.File).Name())
		switch ext {
		case "", ".jpg", ".jpeg":
			err = jpeg.Encode(dst, img, &jpeg.Options{Quality: 100})
		case ".png":
			err = png.Encode(dst, img)
		default:
			err = errors.New("unsupported image format")
		}
	default:
		err = jpeg.Encode(dst, img, &jpeg.Options{Quality: 100})
	}
	return err
}

/*
// inSlice checks if the item exists in the slice.
func inSlice(item string, slice []string) bool {
	for _, it := range slice {
		if it == item {
			return true
		}
	}
	return false
}
*/

// drawEyeDetectionMarker is a helper function to draw the detection marks
func drawEyeDetectionMarker(ctx *gg.Context, x, y, r float64, c color.RGBA, markDet bool) {
	ctx.DrawArc(x, y, r*0.15, 0, 2*math.Pi)
	ctx.SetFillStyle(gg.NewSolidPattern(c))
	ctx.Fill()

	if markDet {
		ctx.DrawRectangle(x-(r*1.5), y-(r*1.5), r*3, r*3)
		ctx.SetLineWidth(2.0)
		ctx.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{R: 255, G: 255, B: 0, A: 255}))
		ctx.Stroke()
	}
}

// coord holds the detection coordinates
type coord struct {
	Row   int `json:"x,omitempty"`
	Col   int `json:"y,omitempty"`
	Scale int `json:"size,omitempty"`
}

// detection holds the detection points of the various detection types
type detection struct {
	FacePoints     coord   `json:"face,omitempty"`
	EyePoints      []coord `json:"eyes,omitempty"`
	LandmarkPoints []coord `json:"landmark_points,omitempty"`
}
