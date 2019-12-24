package cmd

import (
	"github.com/daisaru11/cloud-secrets/cmd/controller"
	"github.com/daisaru11/cloud-secrets/cmd/exec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud-secrets",
		Short: "Decodes and mounts secrets provided from cloud secret stores",
		Long:  `Decodes and mounts secrets provided from cloud secret stores`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			debug, _ := cmd.Flags().GetBool("debug")
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
				logrus.Debugln("enable debug logging.")
			}
		},
	}

	cmd.PersistentFlags().Bool("debug", false, "Enable debug output in logs")

	cmd.AddCommand(controller.NewControllerCommand())
	cmd.AddCommand(exec.NewExecCommand())

	return cmd
}
