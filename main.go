package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"vault-test/envvar"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

type Config struct {
	Imam string `mapstructure:"imam"`
}

func LoadConfig(path io.Reader) (config Config, err error) {
	viper.SetConfigType("json")
	err = viper.ReadConfig(path)
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

func main() {
	if err := envvar.Load(".env"); err != nil {
		log.Fatalln("Couldn't load configuration", err)
	}

	token, err := UserPassLogin()
	if err != nil {
		panic(err)
	}
	client, err := api.NewClient(&api.Config{Address: os.Getenv("VAULT_ADDRESS"), HttpClient: httpClient})
	if err != nil {
		panic(err)
	}

	client.SetToken(token)
	data, err := client.Logical().Read(os.Getenv("VAULT_PATH"))
	if err != nil {
		panic(err)
	}

	b, _ := json.Marshal(data.Data["data"])
	read := bytes.NewReader(b)
	config, _ := LoadConfig(read)
	fmt.Println(config.Imam)
}

func UserPassLogin() (string, error) {
	// create a vault client
	client, err := api.NewClient(&api.Config{Address: os.Getenv("VAULT_ADDRESS"), HttpClient: httpClient})
	if err != nil {
		return "", err
	}

	// to pass the password
	options := map[string]interface{}{
		"password": os.Getenv("VAULT_PASSWORD"),
	}
	path := fmt.Sprintf("auth/userpass/login/%s", os.Getenv("VAULT_USERNAME"))

	// PUT call to get a token
	secret, err := client.Logical().Write(path, options)
	if err != nil {
		return "", err
	}

	token := secret.Auth.ClientToken
	return token, nil
}
