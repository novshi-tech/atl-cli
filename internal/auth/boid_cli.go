package auth

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// boidCLINamespace is the boid secret store namespace atl's credentials are
// kept under (`boid secret ... --namespace atl`).
const boidCLINamespace = "atl"

// BoidCLIStore implements CredentialStore by shelling out to `boid secret`.
//
// This is meant for atl running as a boid host_command, inside the boid
// daemon container: the daemon's own UNIX socket (which `boid secret ...`
// dials by default via XDG_RUNTIME_DIR) is a trusted transport with no
// auth layer at all for anything running inside that same container, so no
// credential handoff of its own is needed here — the `boid` subprocess just
// talks to the daemon it's colocated with. Neither the OS keyring nor
// pass(1) work in that environment (no desktop session, no GPG key), which
// is why this backend exists as a third option rather than trying to make
// either of those fit.
//
// Unlike KeyringStore/PassStore, this backend is never auto-selected by
// NewStore — see its doc comment for why an explicit opt-in is required.
type BoidCLIStore struct{}

// NewBoidCLIStore creates a BoidCLIStore if the boid binary is available on
// PATH.
func NewBoidCLIStore() (*BoidCLIStore, error) {
	if _, err := exec.LookPath("boid"); err != nil {
		return nil, fmt.Errorf("boid not found in PATH: %w", err)
	}
	return &BoidCLIStore{}, nil
}

func (b *BoidCLIStore) Set(key, value string) error {
	var stderr bytes.Buffer
	cmd := exec.Command("boid", "secret", "set", key, "--namespace", boidCLINamespace)
	cmd.Stdin = strings.NewReader(value)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("boid secret set %s: %w: %s", key, err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func (b *BoidCLIStore) Get(key string) (string, error) {
	var stderr bytes.Buffer
	cmd := exec.Command("boid", "secret", "get", key, "--namespace", boidCLINamespace)
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("boid secret get %s: %w: %s", key, err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(string(out)), nil
}

func (b *BoidCLIStore) Delete(key string) error {
	var stderr bytes.Buffer
	cmd := exec.Command("boid", "secret", "delete", key, "--namespace", boidCLINamespace)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("boid secret delete %s: %w: %s", key, err, strings.TrimSpace(stderr.String()))
	}
	return nil
}
