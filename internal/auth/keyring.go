package auth

import "github.com/zalando/go-keyring"

const keyringService = "jira-cli"

// KeyringStore implements CredentialStore using the OS keyring.
type KeyringStore struct{}

// Probe tests whether the keyring is available by writing and deleting a test entry.
func (k *KeyringStore) Probe() error {
	const probeUser = "__probe__"
	if err := keyring.Set(keyringService, probeUser, "test"); err != nil {
		return err
	}
	return keyring.Delete(keyringService, probeUser)
}

func (k *KeyringStore) Set(key, value string) error {
	return keyring.Set(keyringService, key, value)
}

func (k *KeyringStore) Get(key string) (string, error) {
	return keyring.Get(keyringService, key)
}

func (k *KeyringStore) Delete(key string) error {
	return keyring.Delete(keyringService, key)
}
