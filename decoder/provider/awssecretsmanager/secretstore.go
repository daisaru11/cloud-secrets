package awssecretsmanager

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/daisaru11/cloud-secrets/decoder/variable"
)

type SecretStore struct {
	sess *session.Session
}

func NewSecretStore() *SecretStore {
	return &SecretStore{}
}

func (s *SecretStore) GetValue(sv *variable.SecretVariable) (string, error) {
	err := s.initSessionIfNotExist()
	if err != nil {
		return "", err
	}

	smsvc := secretsmanager.New(s.sess)

	id := sv.Identifier
	if !strings.HasPrefix(id, "/") {
		id = "/" + id
	}

	output, err := smsvc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(id),
	})
	if err != nil {
		return "", err
	}

	return aws.StringValue(output.SecretString), nil
}

func (s *SecretStore) initSessionIfNotExist() error {
	if s.sess == nil {
		sess, err := session.NewSession()
		if err != nil {
			return err
		}

		s.sess = sess
	}

	return nil
}
