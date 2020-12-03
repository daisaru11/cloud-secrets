package env

import (
	"os"
	"strings"

	"github.com/daisaru11/cloud-secrets/decoder"
)

func ReplaceEnvironmentVariable(name string) error {
	decoder := decoder.NewDecoder()

	vars := map[string]string{
		name: os.Getenv(name),
	}

	decodedVars, err := decoder.DecodeVariables(vars)
	if err != nil {
		return err
	}

	for name, value := range decodedVars {
		os.Setenv(name, value)
	}

	return nil
}

func ReplaceEnvironmentVariables() error {
	vars := getEnvironmentVariables()

	decoder := decoder.NewDecoder()

	decodedVars, err := decoder.DecodeVariables(vars)
	if err != nil {
		return err
	}

	for name, value := range decodedVars {
		os.Setenv(name, value)
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
