package vault

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/api"
)

type ConfigVault struct {
	Address  string
	Token    string
	Username string
	Password string
	RoleID   string
	SecretID string
	Path     string
}

type Configuration struct {
	vaultClient *api.Client
}

func NewVaultProvider(typeAuth string, config *ConfigVault) (interface{}, error) {
	provider := ClientApi(config.Address)

	switch typeAuth {
	case "token":
		results, err := provider.ReadPath(config.Token, config.Path)
		if err != nil {
			return nil, errors.New("couldn't load provider")
		}
		return results, nil
	case "userpass":
		token, err := provider.UserPassLogin(config.Username, config.Password)
		if err != nil {
			return nil, fmt.Errorf("unable to log in with userpass: %w", err)
		}
		results, err := provider.ReadPath(token, config.Path)
		if err != nil {
			return nil, errors.New("couldn't load provider")
		}
		return results, nil
	case "approle":
		token, err := provider.AppRoleLogin(config.RoleID, config.SecretID)
		if err != nil {
			return nil, fmt.Errorf("unable to log in with approle: %w", err)
		}
		results, err := provider.ReadPath(token, config.Path)
		if err != nil {
			return nil, errors.New("couldn't load provider")
		}
		return results, nil
	default:
		return nil, errors.New("couldn't load provider")
	}
}

// New Vault Token
func (c Configuration) ReadPath(token, path string) (interface{}, error) {

	c.vaultClient.SetToken(token)

	secret, err := c.vaultClient.Logical().Read(path)
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

// Username Password Login
func (c Configuration) UserPassLogin(username, password string) (string, error) {
	// to pass the password
	options := map[string]interface{}{
		"password": password,
	}
	path := fmt.Sprintf("auth/userpass/login/%s", username)

	// PUT call to get a token
	secret, err := c.vaultClient.Logical().Write(path, options)
	if err != nil {
		return "", err
	}

	token := secret.Auth.ClientToken
	return token, nil
}

// AppRole Login
func (c Configuration) AppRoleLogin(roleID, secretID string) (string, error) {
	// to pass the password
	options := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	// PUT call to get a token
	secret, err := c.vaultClient.Logical().Write("auth/approle/login", options)
	if err != nil {
		return "", err
	}

	token := secret.Auth.ClientToken
	return token, nil
}

func ClientApi(address string) *Configuration {
	apiConfig := &api.Config{
		Address: address,
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	client, err := api.NewClient(apiConfig)
	if err != nil {
		return nil
	}
	return &Configuration{client}
}
