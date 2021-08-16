package tools

import (
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/bimg"
	"github.com/sirupsen/logrus"
)

func Blur(path string, suffix string, ratio float64, quality int) (err error) {
	input, _ := os.Open(path)
	defer input.Close()

	dir, filenameWithExt := filepath.Split(path)
	filename := strings.TrimSuffix(filenameWithExt, filepath.Ext(path))

	buf, err := bimg.Read(path)
	if err != nil {
		logrus.Error(err)
		return
	}

	img := bimg.NewImage(buf)
	size, _ := img.Size()
	newImage, err := img.Process(bimg.Options{
		Width:   int(math.Ceil(float64(size.Width) * ratio)),
		Height:  int(math.Ceil(float64(size.Height) * ratio)),
		Quality: quality,
		GaussianBlur: bimg.GaussianBlur{
			Sigma:   10,
			MinAmpl: 10,
		},
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	err = bimg.Write(filepath.Join(dir, filename+"-"+suffix+".jpg"), newImage)
	if err != nil {
		logrus.Error(err)
		return
	}

	return
}
