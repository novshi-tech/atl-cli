package auth

import (
	"fmt"
	"strings"
)

// CredentialStore defines an interface for storing and retrieving credentials.
type CredentialStore interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

// SiteCredentials holds authentication information for a Jira site.
type SiteCredentials struct {
	JiraURL  string
	Email    string
	APIToken string
}

// NewStore returns a CredentialStore, trying keyring first, then falling back to pass.
func NewStore() (CredentialStore, error) {
	ks := &KeyringStore{}
	if err := ks.Probe(); err == nil {
		return ks, nil
	}
	ps, err := NewPassStore()
	if err != nil {
		return nil, fmt.Errorf("no credential backend available: keyring unavailable, pass not found")
	}
	return ps, nil
}

// SaveSite stores credentials for the given site alias.
func SaveSite(store CredentialStore, alias string, creds SiteCredentials) error {
	if err := store.Set(alias+"/jira-url", creds.JiraURL); err != nil {
		return fmt.Errorf("saving jira-url: %w", err)
	}
	if err := store.Set(alias+"/email", creds.Email); err != nil {
		return fmt.Errorf("saving email: %w", err)
	}
	if err := store.Set(alias+"/api-token", creds.APIToken); err != nil {
		return fmt.Errorf("saving api-token: %w", err)
	}
	return nil
}

// LoadSite loads credentials for the given site alias.
func LoadSite(store CredentialStore, alias string) (SiteCredentials, error) {
	url, err := store.Get(alias + "/jira-url")
	if err != nil {
		return SiteCredentials{}, fmt.Errorf("loading jira-url for site %q: %w", alias, err)
	}
	email, err := store.Get(alias + "/email")
	if err != nil {
		return SiteCredentials{}, fmt.Errorf("loading email for site %q: %w", alias, err)
	}
	token, err := store.Get(alias + "/api-token")
	if err != nil {
		return SiteCredentials{}, fmt.Errorf("loading api-token for site %q: %w", alias, err)
	}
	return SiteCredentials{JiraURL: url, Email: email, APIToken: token}, nil
}

// GetDefaultSite returns the default site alias.
func GetDefaultSite(store CredentialStore) (string, error) {
	return store.Get("default-site")
}

// SetDefaultSite sets the default site alias.
func SetDefaultSite(store CredentialStore, alias string) error {
	return store.Set("default-site", alias)
}

// ListSites returns all configured site aliases.
func ListSites(store CredentialStore) ([]string, error) {
	raw, err := store.Get("sites")
	if err != nil {
		return nil, nil // no sites configured yet
	}
	if raw == "" {
		return nil, nil
	}
	return strings.Split(raw, ","), nil
}

// AddSiteToList adds a site alias to the sites list if not already present.
func AddSiteToList(store CredentialStore, alias string) error {
	sites, err := ListSites(store)
	if err != nil {
		return err
	}
	for _, s := range sites {
		if s == alias {
			return nil // already present
		}
	}
	sites = append(sites, alias)
	return store.Set("sites", strings.Join(sites, ","))
}
