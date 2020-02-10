// Package main implements the go-service command.
package main

import (
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/golang/glog"
	"github.com/spf13/viper"

	"github.com/jacobgarcia/settify/fixer"
	"github.com/jacobgarcia/settify/server"
)

func main() {
	confDirFlag := flag.String("conf-dir", "conf", "custom conf dir folder")
	flag.Parse()

	glog.Info("starting..")

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	viper := viper.New()
	viper.SetConfigName("settify")
	viper.AddConfigPath(*confDirFlag)

	// Use env for values if declared
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		glog.Exitf("Fatal error config file: %s", err)
	}

	fixerClient := fixer.New(viper.GetString("fixer.URL"), viper.GetString("fixer.key"))

	handler := server.MakeHTTPHandler(fixerClient, logger)
	port := viper.GetString("port")

	if err != nil {
		glog.Exitf("Error starting the server: %s", err)
	}

	glog.Info("serving on port %s", port)
	err = http.ListenAndServe(":"+port, handler)
	glog.Exitf("Server stopped: %s", err)
}
