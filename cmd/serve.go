package cmd

import (
	"net/http"

	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/pkg/api"
	"github.com/last-second/services/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs the server that provides an interface to a local instance of the database",
	Long:  "Runs the server that provides an interface to a local instance of the database",
	Run:   runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) {
	config.Init()
	logrus.Info("starting serve")

	if err := http.ListenAndServe(":8000", api.New()); err != nil {
		logrus.Fatal(trace.New("ErrorApiRuntime").Trace(err))
	}
}
