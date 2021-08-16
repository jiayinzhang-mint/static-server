package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/h2non/bimg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// BuildEnv build mode (dev or prod)
var BuildEnv string

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true).PathPrefix("/api").Subrouter()
	r.HandleFunc("/image", func(rw http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		path := params["path"]

		// Invalid path
		if strings.Contains(path, "..") {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		fullPath := filepath.Join(viper.GetString("file.upload"), path)
		fileInfo, e := os.Stat(fullPath)

		// File not exist
		if os.IsNotExist(e) {
			logrus.Error(path, " does not exists.")
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		// Path is a dir
		if fileInfo.IsDir() {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check if file exists and open
		Openfile, err := os.Open(fullPath)
		if err != nil {
			// File not found, send 404
			http.Error(rw, "File not found.", 404)
		}
		defer Openfile.Close() // Close after function return

		http.ServeFile(rw, r, fullPath)
	}).Queries("path", "{path}").Methods("GET")

	log.Fatal(http.ListenAndServe(":9099", r))
}

func main() {
	logrus.Info("Server listening at 9099")
	logrus.Info("Build Env: " + BuildEnv)

	// Get config from json
	if BuildEnv == "prod" {
		viper.SetConfigName("config.prod")
	} else {
		viper.SetConfigName("config.test")
	}

	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("config file error")
	}

	buffer, err := bimg.Read("./upload/1fab7ccc-0b9a-4769-a567-1d09a6484276/biu.jpg")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	_, err = bimg.NewImage(buffer).Resize(800, 600)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	logrus.Info("press ctrl+c to exit")
	handleRequests()
}
