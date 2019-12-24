package provider

import (
	"github.com/daisaru11/cloud-secrets/decoder/variable"
)

type SecretStore interface {
	GetValue(sv *variable.SecretVariable) (string, error)
}
