package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"static-server/tools"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// BuildEnv build mode (dev or prod)
var BuildEnv string

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true).PathPrefix("/api").Subrouter()
	r.HandleFunc("/image", func(rw http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		path := params.Get("path")
		filter := params.Get("filter")

		// Filter validity
		if filter != "" {
			filters_available := []string{"blur"}
			valid := func() bool {
				for _, v := range filters_available {
					if filter == v {
						return true
					}
				}
				return false
			}()
			if !valid {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		// Path validity
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

		if filter == "blur" {
			blurredFilePath := strings.TrimSuffix(fullPath, filepath.Ext(path)) + "-blur" + filepath.Ext(path)
			if _, notExist := os.Stat(blurredFilePath); os.IsNotExist(notExist) {
				err := tools.Blur(fullPath, "blur", 0.5, 90)
				if err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			// Check if file exists and open
			Openfile, err := os.Open(blurredFilePath)
			if err != nil {
				http.Error(rw, "File not found.", 404)
			}
			defer Openfile.Close()

			http.ServeFile(rw, r, blurredFilePath)
			return
		}

		// Check if file exists and open
		Openfile, err := os.Open(fullPath)
		if err != nil {
			http.Error(rw, "File not found.", 404)
		}
		defer Openfile.Close()

		http.ServeFile(rw, r, fullPath)
	}).Methods("GET")

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

	logrus.Info("press ctrl+c to exit")
	handleRequests()
}
