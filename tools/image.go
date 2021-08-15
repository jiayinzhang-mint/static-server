package tools

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/image/draw"
)

func CompressWithWidth(path string, suffix string, width int, quality int) (err error) {
	input, _ := os.Open(path)
	defer input.Close()

	dir, filenameWithExt := filepath.Split(path)
	filename := strings.TrimSuffix(filenameWithExt, filepath.Ext(path))
	output, err := os.Create(filepath.Join(dir, filename+"-"+suffix+".jpg"))
	if err != nil {
		logrus.Error(err)
		return
	}

	defer output.Close()

	// Decode the image (from PNG to image.Image):
	src, err := jpeg.Decode(input)
	if err != nil {
		logrus.Error(err)
		return
	}

	height := width * src.Bounds().Max.Y / src.Bounds().Max.X

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	// Resize:
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	// Encode to `output`:
	err = jpeg.Encode(output, dst, &jpeg.Options{Quality: quality})
	if err != nil {
		logrus.Error(err)
		return
	}

	return
}
