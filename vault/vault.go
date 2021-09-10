package vault

import (
	"errors"
	"fmt"

	"github.com/hashicorp/vault/api"
)

type ConfigVault struct {
	Address  string
	Token    string
	Username string
	Password string
	RoleID   string
	SecretID string
}

type Configuration struct {
	client *api.Logical
}

func NewVaultProvider(typeAuth string, config *ConfigVault) (*Configuration, error) {
	vaultApi := New(config.Address)

	switch typeAuth {
	case "token":
		vaultApi.client.SetToken(config.Token)
		return &Configuration{
			client: vaultApi.client.Logical(),
		}, nil
	case "userpass":
		token, err := vaultApi.UserPassLogin(config.Username, config.Password)
		if err != nil {
			return nil, fmt.Errorf("unable to log in with userpass: %w", err)
		}
		vaultApi.client.SetToken(token)
		return &Configuration{
			client: vaultApi.client.Logical(),
		}, nil
	case "approle":
		token, err := vaultApi.AppRoleLogin(config.RoleID, config.SecretID)
		if err != nil {
			return nil, fmt.Errorf("unable to log in with approle: %w", err)
		}
		vaultApi.client.SetToken(token)
		return &Configuration{
			client: vaultApi.client.Logical(),
		}, nil
	default:
		return nil, errors.New("couldn't load provider")
	}
}

func (c Configuration) Get(path string) (interface{}, error) {

	secret, err := c.client.Read(path)
	if err != nil {
		return "", fmt.Errorf("reading: %w", err)
	}

	if secret == nil {
		return "", errors.New("secret not found")
	}

	data, ok := secret.Data["data"]

	if !ok {
		return "", errors.New("invalid data in secret")
	}
	return data, nil
}
