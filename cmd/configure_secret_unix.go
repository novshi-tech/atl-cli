//go:build !windows

package cmd

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

func readSecretDetectPaste() (string, error) {
	fd := syscall.Stdin

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// Fallback: use standard ReadPassword; indicator shows only after Enter.
		b, e := term.ReadPassword(fd)
		if len(b) > 0 {
			fmt.Println(" (received)")
		} else {
			fmt.Println()
		}
		return string(b), e
	}
	defer term.Restore(fd, oldState)

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
