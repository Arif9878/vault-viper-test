package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	vaultRoleID := os.Getenv("VAULT_ROLE_ID")
	vaultSecretID := os.Getenv("VAULT_SECRET_ID")
	vaultAddress := os.Getenv("VAULT_ADDRESS")

	data, err := vault.NewVaultProvider("approle", &vault.ConfigVault{
		RoleID:   vaultRoleID,
		SecretID: vaultSecretID,
		Address:  vaultAddress,
		Path:     vaultPath,
	})
	if err != nil {
		fmt.Println(err)
	}
	b, _ := json.Marshal(data)
	read := bytes.NewReader(b)
	config, _ := config.LoadConfig(read)
	fmt.Println(config)
}
