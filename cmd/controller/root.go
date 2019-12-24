package controller

import (
	"net/http"

	"github.com/daisaru11/cloud-secrets/webhook/mutating"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewControllerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "controller",
		Short: "cloud-secrets admission controller",
		RunE:  runControllerCommand,
	}

	return cmd
}

func runControllerCommand(cmd *cobra.Command, args []string) error {
	addr := ":8080"

	logrus.Debugln("Starting webhook server")
	logrus.Debugf("Listening on %s", addr)

	certFile := "/certs/tls.crt"
	keyFile := "/certs/tls.key"

	mutatingHandler := mutating.NewHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", mutatingHandler.Handle)
	// mux.HandleFunc("/health/ready", readyHandler)
	err := http.ListenAndServeTLS(":8080", certFile, keyFile, mux)
	if err != nil {
		return err
	}

	return nil
}
