package vault

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/api"
)

type Client struct {
	client *api.Client
}

func New(address string) *Client {
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
	return &Client{client}
}

// Username Password Login
func (c Client) UserPassLogin(username, password string) (string, error) {
	// to pass the password
	options := map[string]interface{}{
		"password": password,
	}
	path := fmt.Sprintf("auth/userpass/login/%s", username)

	// PUT call to get a token
	secret, err := c.client.Logical().Write(path, options)
	if err != nil {
		return "", err
	}

	token := secret.Auth.ClientToken
	return token, nil
}

// AppRole Login
func (c Client) AppRoleLogin(roleID, secretID string) (string, error) {
	// to pass the password
	options := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	// PUT call to get a token
	secret, err := c.client.Logical().Write("auth/approle/login", options)
	if err != nil {
		return "", err
	}

	token := secret.Auth.ClientToken
	return token, nil
}
