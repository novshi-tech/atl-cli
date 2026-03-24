package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

const (
	moduleProxyBase = "https://proxy.golang.org"
	modulePath      = "github.com/novshi-tech/atl-cli"
	checkDateFile   = "last_update_check"
)

type proxyLatestInfo struct {
	Version string `json:"Version"`
}

func atlCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".atl")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return dir, nil
}

func readLastCheckDate() string {
	dir, err := atlCacheDir()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(dir, checkDateFile))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func saveLastCheckDate(date string) {
	dir, err := atlCacheDir()
	if err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(dir, checkDateFile), []byte(date), 0600)
}

func currentBuildVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return info.Main.Version
}

func fetchLatestVersion() (string, error) {
	url := fmt.Sprintf("%s/%s/@latest", moduleProxyBase, modulePath)
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("proxy returned HTTP %d", resp.StatusCode)
	}
	var info proxyLatestInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", err
	}
	return info.Version, nil
}

// checkForUpdate checks once per day if a newer version is available and
// notifies the user if one is found.
func checkForUpdate() {
	today := time.Now().Format("2006-01-02")

	if readLastCheckDate() == today {
		return
	}

	// Mark the check as done for today before anything else so that network
	// errors don't cause repeated slow attempts.
	saveLastCheckDate(today)

	current := currentBuildVersion()
	if current == "" || current == "(devel)" {
		// Dev/local builds — skip check.
		return
	}

	latest, err := fetchLatestVersion()
	if err != nil || latest == "" || latest == current {
		return
	}

	fmt.Fprintf(os.Stderr, "atl: new version available (%s → %s)\n", current, latest)
	fmt.Fprintf(os.Stderr, "atl: run `go install %s` to update.\n", modulePath+"/cmd/atl@latest")
}
