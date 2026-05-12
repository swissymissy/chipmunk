package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// prompt for entering password twice
func PromptPasswordTwice() string {
	fmt.Print("Enter new password: ")
	pw1, err := term.ReadPassword(int(os.Stdin.Fd())) // read input from term
	fmt.Println()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading password: %v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(string(pw1)) == "" {
		fmt.Fprintln(os.Stderr, "password cannot be empty")
		os.Exit(1)
	}

	fmt.Print("Confirm password: ")
	pw2, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading confirmation password: %v\n", err)
		os.Exit(1)
	}

	// compare 2 passwords
	if string(pw1) != string(pw2) {
		fmt.Fprintln(os.Stderr, "passwords do not match")
		os.Exit(1)
	}
	return string(pw1)
}
