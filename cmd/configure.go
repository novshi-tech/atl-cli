package cmd

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/novshi-tech/atl-cli/internal/auth"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure authentication for an Atlassian site",
	RunE:  runConfigure,
}

func init() {
	configureCmd.Flags().StringP("site", "s", "", "Site alias name (required)")
	configureCmd.MarkFlagRequired("site")
	configureCmd.Flags().Bool("default", false, "Set this site as the default")
	rootCmd.AddCommand(configureCmd)
}

func runConfigure(cmd *cobra.Command, args []string) error {
	alias, _ := cmd.Flags().GetString("site")
	setDefault, _ := cmd.Flags().GetBool("default")

	store, err := auth.NewStore()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)

	// Load existing credentials so Enter preserves current values
	existing, _ := auth.LoadSite(store, alias)

	fmt.Printf("Configuring site %q\n", alias)

	siteURL := promptText(reader, "Site URL", existing.BaseURL)
	email := promptText(reader, "Email", existing.Email)
	apiToken := promptSecret("API Token", existing.APIToken)
	bbAPIToken := promptSecret("Bitbucket API Token", existing.BBAPIToken)

	creds := auth.SiteCredentials{
		BaseURL:    siteURL,
		Email:      email,
		APIToken:   apiToken,
		BBAPIToken: bbAPIToken,
	}

	if err := auth.SaveSite(store, alias, creds); err != nil {
		return fmt.Errorf("saving credentials: %w", err)
	}

	if err := auth.AddSiteToList(store, alias); err != nil {
		return fmt.Errorf("updating site list: %w", err)
	}

	// Auto-set as default if this is the first site or --default is specified
	if setDefault {
		if err := auth.SetDefaultSite(store, alias); err != nil {
			return fmt.Errorf("setting default site: %w", err)
		}
		fmt.Printf("Site %q set as default.\n", alias)
	} else {
		sites, _ := auth.ListSites(store)
		if len(sites) == 1 {
			if err := auth.SetDefaultSite(store, alias); err != nil {
				return fmt.Errorf("setting default site: %w", err)
			}
			fmt.Printf("Site %q set as default (first site).\n", alias)
		}
	}

	fmt.Printf("Site %q configured successfully.\n", alias)
	return nil
}

// promptText prompts for a text value. If the user presses Enter, the existing value is kept.
func promptText(reader *bufio.Reader, label, current string) string {
	if current != "" {
		fmt.Printf("%s [%s]: ", label, current)
	} else {
		fmt.Printf("%s: ", label)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return current
	}
	return input
}

// promptSecret prompts for a secret value (input hidden). If the user presses Enter, the existing value is kept.
// The indicator "(parsed from clipboard)" or "(received)" appears immediately upon input, before Enter is pressed.
func promptSecret(label, current string) string {
	if current != "" {
		fmt.Printf("%s [****]: ", label)
	} else {
		fmt.Printf("%s (Enter to skip): ", label)
	}
	secret, err := readSecretDetectPaste()
	if err != nil || len(secret) == 0 {
		return current
	}
	return secret
}

// readSecretDetectPaste reads a secret in raw/no-echo mode and prints an indicator immediately
// upon detecting input — "(parsed from clipboard)" for Cmd+V / Ctrl+V paste (detected via
// terminal bracketed paste mode: ESC[200~ start, ESC[201~ end), or "(received)" for typed input.
// The indicator appears right when input arrives, not after the user presses Enter.
func readSecretDetectPaste() (string, error) {
	// On Windows, always use fallback since MakeRaw and bracketed paste don't work reliably
	if runtime.GOOS == "windows" {
		b, e := term.ReadPassword(int(syscall.Stdin))
		if len(b) > 0 {
			fmt.Println(" (received)")
		} else {
			fmt.Println()
		}
		return string(b), e
	}

	fd := syscall.Handle(syscall.Stdin)

	oldState, err := term.MakeRaw(int(fd))
	if err != nil {
		// Fallback: use standard ReadPassword; indicator shows only after Enter.
		b, e := term.ReadPassword(int(fd))
		if len(b) > 0 {
			fmt.Println(" (received)")
		} else {
			fmt.Println()
		}
		return string(b), e
	}
	defer term.Restore(int(fd), oldState)

	// Enable bracketed paste mode; disable it on return.
	fmt.Fprint(os.Stdout, "\x1b[?2004h")
	defer fmt.Fprint(os.Stdout, "\x1b[?2004l")

	// State machine to parse paste escape sequences:
	//   paste start: ESC [ 2 0 0 ~  (\x1b[200~)
	//   paste end:   ESC [ 2 0 1 ~  (\x1b[201~)
	const (
		stNormal   = iota
		stEsc      // saw \x1b
		stEscBr    // saw \x1b[
		stEscBr2   // saw \x1b[2
		stEscBr20  // saw \x1b[20
		stEscBr200 // saw \x1b[200 — next '~' completes paste start
		stEscBr201 // saw \x1b[201 — next '~' completes paste end
	)

	var result []byte
	indicatorShown := false
	state := stNormal
	oneByte := make([]byte, 1)

	// showIndicator prints msg immediately on the same prompt line, without a newline.
	// The newline is deferred until the user presses Enter.
	showIndicator := func(msg string) {
		if !indicatorShown {
			fmt.Print(msg)
			indicatorShown = true
		}
	}

	for {
		n, readErr := syscall.Read(fd, oneByte)
		if readErr != nil || n == 0 {
			break
		}
		b := oneByte[0]

		switch state {
		case stNormal:
			switch b {
			case '\r', '\n':
				fmt.Print("\r\n")
				return string(result), nil
			case 3, 4: // Ctrl+C, Ctrl+D
				fmt.Print("\r\n")
				return "", fmt.Errorf("interrupted")
			case 0x7f, 0x08: // Backspace / Delete
				if len(result) > 0 {
					result = result[:len(result)-1]
				}
			case 0x1b: // ESC — possible start of paste sequence
				state = stEsc
			default:
				showIndicator(" (received)")
				result = append(result, b)
			}
		case stEsc:
			if b == '[' {
				state = stEscBr
			} else {
				state = stNormal
			}
		case stEscBr:
			if b == '2' {
				state = stEscBr2
			} else {
				state = stNormal
			}
		case stEscBr2:
			if b == '0' {
				state = stEscBr20
			} else {
				state = stNormal
			}
		case stEscBr20:
			switch b {
			case '0':
				state = stEscBr200
			case '1':
				state = stEscBr201
			default:
				state = stNormal
			}
		case stEscBr200:
			if b == '~' {
				showIndicator(" (parsed from clipboard)") // shown immediately on paste start
			}
			state = stNormal
		case stEscBr201:
			state = stNormal // paste end sequence consumed
		}
	}

	fmt.Print("\r\n")
	return string(result), nil
}
