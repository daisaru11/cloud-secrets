package decoder

import (
	"fmt"

	"github.com/daisaru11/cloud-secrets/decoder/provider"
	"github.com/daisaru11/cloud-secrets/decoder/provider/awsparameterstore"
	"github.com/daisaru11/cloud-secrets/decoder/provider/awssecretsmanager"
	"github.com/daisaru11/cloud-secrets/decoder/provider/nop"
	"github.com/daisaru11/cloud-secrets/decoder/variable"
)

type Decoder struct {
	secretStores map[string]provider.SecretStore
}

func NewDecoder() *Decoder {
	return &Decoder{
		secretStores: map[string]provider.SecretStore{},
	}
}

func (e *Decoder) DecodeVariables(vars map[string]string) (map[string]string, error) {
	secretVars := variable.ExtractSecretVariables(vars)
	decodedVars := map[string]string{}

	for name, sv := range secretVars {
		secretVar := sv

		secretStore, err := e.getSecretStoreByProviderName(secretVar.Provider)
		if err != nil {
			return nil, err
		}

		value, err := secretStore.GetValue(&secretVar)
		if err != nil {
			return nil, err
		}

		decodedVars[name] = value
	}

	return decodedVars, nil
}

func (e *Decoder) getSecretStoreByProviderName(provider string) (provider.SecretStore, error) {
	ss, ok := e.secretStores[provider]
	if ok {
		return ss, nil
	}

	switch provider {
	case "aws-parameter-store":
		ss = awsparameterstore.NewSecretStore()
		e.secretStores[provider] = ss

		return ss, nil
	case "aws-secrets-manager":
		ss = awssecretsmanager.NewSecretStore()
		e.secretStores[provider] = ss

		return ss, nil
	case "nop":
		ss = nop.NewSecretStore()
		e.secretStores[provider] = ss

		return ss, nil

	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}
