package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/daisaru11/cloud-secrets/decoder"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewExecCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec [command]",
		Short: "Execute command with decrypted variables",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runExecCommand,
	}

	cmd.Flags().SetInterspersed(false)

	return cmd
}

// nolint: gosec
func runExecCommand(cmd *cobra.Command, args []string) error {
	logrus.Debugln("Starting Exec command")
	logrus.Debugf("With the args: %s", strings.Join(args, ", "))

	vars := getEnvironmentVariables()

	decoder := decoder.NewDecoder()

	decodedVars, err := decoder.DecodeVariables(vars)
	if err != nil {
		return err
	}

	command := exec.Command(args[0])
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = append(os.Environ(), formatEnviron(decodedVars)...)
	command.Args = append(command.Args, args[1:]...)

	err = command.Run()
	if err != nil {
		return fmt.Errorf("failed to run docker cli: %w", err)
	}

	return nil
}

func getEnvironmentVariables() map[string]string {
	envs := os.Environ()
	varsmap := map[string]string{}

	for _, v := range envs {
		splitted := strings.SplitN(v, "=", 2)
		key := splitted[0]
		val := splitted[1]

		varsmap[key] = val
	}

	return varsmap
}

func formatEnviron(envs map[string]string) []string {
	environ := []string{}
	for name, value := range envs {
		environ = append(environ, fmt.Sprintf("%s=%s", name, value))
	}

	return environ
}
