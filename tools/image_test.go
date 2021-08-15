package tools

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config.test")
	viper.AddConfigPath("..")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal(err)
	}
}

func TestCompress(t *testing.T) {
	Init()
	err := CompressWithWidth("../upload/1fab7ccc-0b9a-4769-a567-1d09a6484276/biu.jpg", "lq", 900, 90)
	if err != nil {
		logrus.Fatal(err)
	}
}
