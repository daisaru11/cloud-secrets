package variable

import (
	"strings"
)

type SecretVariable struct {
	Provider   string
	Identifier string
}

const (
	SecretVariableScheme = "cloud-secrets://"
)

func ExtractSecretVariables(variables map[string]string) map[string]SecretVariable {
	ret := map[string]SecretVariable{}

	for k, v := range variables {
		if !strings.HasPrefix(v, SecretVariableScheme) {
			continue
		}

		trimmed := strings.TrimPrefix(v, SecretVariableScheme)
		splitted := strings.SplitN(trimmed, "/", 2)

		if len(splitted) != 2 {
			continue
		}

		provider := splitted[0]
		identifier := splitted[1]

		ret[k] = SecretVariable{
			Provider:   provider,
			Identifier: identifier,
		}
	}

	return ret
}
