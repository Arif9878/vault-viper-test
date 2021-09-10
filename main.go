package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"vault-test/config"
	"vault-test/envvar"
	"vault-test/vault"
)

func main() {
	if err := envvar.Load(".env"); err != nil {
		log.Fatalln("Couldn't load configuration", err)
	}

	vaultPath := os.Getenv("VAULT_PATH")
	// vaultUsername := os.Getenv("VAULT_USERNAME")
	// vaultPassword := os.Getenv("VAULT_PASSWORD")
	vaultAddress := os.Getenv("VAULT_ADDRESS")
	vaultToken := os.Getenv("VAULT_TOKEN")
	// client, _ := vaultclient.NewClient(&vaultclient.VaultConfig{Server: vaultAddress})
	// client.TokenAuth(vaultToken)
	// data, _ := client.GetValue(vaultPath)
	// fmt.Println(data)

	vault, err := vault.NewVaultProvider("token", &vault.ConfigVault{
		Token:   vaultToken,
		Address: vaultAddress,
	})

	if err != nil {
		fmt.Println(err)
	}

	data, err := vault.Get(vaultPath)
	if err != nil {
		fmt.Println(err)
	}

	b, _ := json.Marshal(data)
	read := bytes.NewReader(b)
	config, _ := config.LoadConfig(read)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(config.Imam))
	})

	var address = "localhost:9001"
	fmt.Printf("server started at %s\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
