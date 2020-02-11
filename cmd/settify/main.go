// Package main implements the settify command
package main

import (
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/golang/glog"
	"github.com/spf13/viper"

	"github.com/jacobgarcia/settify/server"
	"github.com/jacobgarcia/settify/spotify"
)

func main() {
	confDirFlag := flag.String("conf-dir", "conf", "custom conf dir folder")
	flag.Parse()

	glog.Info("Starting server...")

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	viper := viper.New()
	viper.SetConfigName("settify")
	viper.AddConfigPath(*confDirFlag)

	// Use env for values if declared
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		glog.Exitf("Fatal error at configuration file: %s", err)
	}

	spotifyClient := spotify.New(viper.GetString("spotify.authURL"), viper.GetString("spotify.URL"), viper.GetString("spotify.id"), viper.GetString("spotify.secret"))

	router := server.CreateRouter(spotifyClient, logger)
	port := viper.GetString("port")

	if err != nil {
		glog.Exitf("Error starting the server: %s", err)
	}

	glog.Info("Serving on port: ", port)
	err = http.ListenAndServe(":"+port, router)
	glog.Exitf("Server stopped: %s", err)
}
