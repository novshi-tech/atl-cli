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

// SiteCredentials holds authentication information for an Atlassian site.
type SiteCredentials struct {
	BaseURL       string
	Email         string
	APIToken      string
	BBAPIToken string // Bitbucket API Token (optional, used instead of APIToken for Bitbucket)
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
	if err := store.Set(alias+"/base-url", creds.BaseURL); err != nil {
		return fmt.Errorf("saving base-url: %w", err)
	}
	if err := store.Set(alias+"/email", creds.Email); err != nil {
		return fmt.Errorf("saving email: %w", err)
	}
	if err := store.Set(alias+"/api-token", creds.APIToken); err != nil {
		return fmt.Errorf("saving api-token: %w", err)
	}
	if creds.BBAPIToken != "" {
		if err := store.Set(alias+"/bb-api-token", creds.BBAPIToken); err != nil {
			return fmt.Errorf("saving bb-api-token: %w", err)
		}
	}
	return nil
}

// LoadSite loads credentials for the given site alias.
func LoadSite(store CredentialStore, alias string) (SiteCredentials, error) {
	url, err := store.Get(alias + "/base-url")
	if err != nil {
		return SiteCredentials{}, fmt.Errorf("loading base-url for site %q: %w", alias, err)
	}
	email, err := store.Get(alias + "/email")
	if err != nil {
		return SiteCredentials{}, fmt.Errorf("loading email for site %q: %w", alias, err)
	}
	token, err := store.Get(alias + "/api-token")
	if err != nil {
		return SiteCredentials{}, fmt.Errorf("loading api-token for site %q: %w", alias, err)
	}
	bbAPIToken, _ := store.Get(alias + "/bb-api-token") // optional, ignore error
	return SiteCredentials{BaseURL: url, Email: email, APIToken: token, BBAPIToken: bbAPIToken}, nil
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

// RemoveSiteFromList removes a site alias from the sites list.
func RemoveSiteFromList(store CredentialStore, alias string) error {
	sites, err := ListSites(store)
	if err != nil {
		return err
	}
	filtered := sites[:0]
	for _, s := range sites {
		if s != alias {
			filtered = append(filtered, s)
		}
	}
	return store.Set("sites", strings.Join(filtered, ","))
}

// DeleteSite removes all credentials for the given site alias.
func DeleteSite(store CredentialStore, alias string) error {
	if err := store.Delete(alias + "/base-url"); err != nil {
		return fmt.Errorf("deleting base-url: %w", err)
	}
	if err := store.Delete(alias + "/email"); err != nil {
		return fmt.Errorf("deleting email: %w", err)
	}
	if err := store.Delete(alias + "/api-token"); err != nil {
		return fmt.Errorf("deleting api-token: %w", err)
	}
	store.Delete(alias + "/bb-api-token") // optional, ignore error
	if err := RemoveSiteFromList(store, alias); err != nil {
		return fmt.Errorf("removing from site list: %w", err)
	}
	defaultSite, _ := GetDefaultSite(store)
	if defaultSite == alias {
		store.Delete("default-site")
	}
	return nil
}
