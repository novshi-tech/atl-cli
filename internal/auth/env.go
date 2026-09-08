package auth

import (
	"fmt"
	"os"
	"strings"
)

// Environment variables that supply site credentials directly, bypassing the
// credential store entirely. Intended for environments where neither the OS
// keyring nor pass(1) is available (CI jobs, sandboxed agent containers) and
// the caller injects credentials through the process environment instead.
const (
	EnvBaseURL     = "ATL_BASE_URL"
	EnvEmail       = "ATL_EMAIL"
	EnvAPIToken    = "ATL_API_TOKEN"
	EnvBBAPIToken  = "ATL_BB_API_TOKEN" // optional
	EnvBBWorkspace = "ATL_BB_WORKSPACE" // optional
)

// LoadFromEnv returns credentials assembled from ATL_BASE_URL / ATL_EMAIL /
// ATL_API_TOKEN (plus the optional Bitbucket variables). ok is false when none
// of the three required variables is set, in which case callers should fall
// back to the credential store. Setting only some of the required variables
// is treated as a configuration error rather than silently ignored, so a typo
// in one name does not quietly route the command to a different site.
func LoadFromEnv() (creds SiteCredentials, ok bool, err error) {
	baseURL := strings.TrimSpace(os.Getenv(EnvBaseURL))
	email := strings.TrimSpace(os.Getenv(EnvEmail))
	token := strings.TrimSpace(os.Getenv(EnvAPIToken))

	if baseURL == "" && email == "" && token == "" {
		return SiteCredentials{}, false, nil
	}

	var missing []string
	if baseURL == "" {
		missing = append(missing, EnvBaseURL)
	}
	if email == "" {
		missing = append(missing, EnvEmail)
	}
	if token == "" {
		missing = append(missing, EnvAPIToken)
	}
	if len(missing) > 0 {
		return SiteCredentials{}, true, fmt.Errorf("environment credentials are incomplete: missing %s (set all of %s, %s, %s or none of them)",
			strings.Join(missing, ", "), EnvBaseURL, EnvEmail, EnvAPIToken)
	}

	return SiteCredentials{
		BaseURL:     baseURL,
		Email:       email,
		APIToken:    token,
		BBAPIToken:  strings.TrimSpace(os.Getenv(EnvBBAPIToken)),
		BBWorkspace: strings.TrimSpace(os.Getenv(EnvBBWorkspace)),
	}, true, nil
}
