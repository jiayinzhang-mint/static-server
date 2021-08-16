package tools

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestCompress(t *testing.T) {

	err := Blur("../upload/1fab7ccc-0b9a-4769-a567-1d09a6484276/biu.jpg", "blur", 0.5, 90)
	if err != nil {
		logrus.Fatal(err)
	}
}
