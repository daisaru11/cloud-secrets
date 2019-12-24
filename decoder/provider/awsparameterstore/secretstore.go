package awsparameterstore

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
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

	ssmsvc := ssm.New(s.sess)

	name := sv.Identifier
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}

	output, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}

	return aws.StringValue(output.Parameter.Value), nil
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
