//go:build windows

package cmd

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

func readSecretDetectPaste() (string, error) {
	b, e := term.ReadPassword(int(syscall.Stdin))
	if len(b) > 0 {
		fmt.Println(" (received)")
	} else {
		fmt.Println()
	}
	return string(b), e
}
