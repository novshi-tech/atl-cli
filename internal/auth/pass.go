package auth

import (
	"fmt"
	"os/exec"
	"strings"
)

const passPrefix = "atl/"

// PassStore implements CredentialStore using pass(1).
type PassStore struct {
	passPath string
}

// NewPassStore creates a PassStore if pass is available on the system.
func NewPassStore() (*PassStore, error) {
	path, err := exec.LookPath("pass")
	if err != nil {
		return nil, fmt.Errorf("pass not found: %w", err)
	}
	return &PassStore{passPath: path}, nil
}

func (p *PassStore) Set(key, value string) error {
	cmd := exec.Command(p.passPath, "insert", "--force", "--echo", passPrefix+key)
	cmd.Stdin = strings.NewReader(value)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pass insert %s: %w", key, err)
	}
	return nil
}

func (p *PassStore) Get(key string) (string, error) {
	cmd := exec.Command(p.passPath, "show", passPrefix+key)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("pass show %s: %w", key, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (p *PassStore) Delete(key string) error {
	cmd := exec.Command(p.passPath, "rm", "-f", passPrefix+key)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pass rm %s: %w: %s", key, err, string(out))
	}
	return nil
}
