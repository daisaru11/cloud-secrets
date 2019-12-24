package nop

import (
	"fmt"

	"github.com/daisaru11/cloud-secrets/decoder/variable"
)

// Returns just the identifier of secret for debugging purpose.
type SecretStore struct {
}

func NewSecretStore() *SecretStore {
	return &SecretStore{}
}

func (s *SecretStore) GetValue(sv *variable.SecretVariable) (string, error) {
	return fmt.Sprintf("cloud-secrets://decoded/%s", sv.Identifier), nil
}
