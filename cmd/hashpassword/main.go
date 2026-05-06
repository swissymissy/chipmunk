package main

import (
	"fmt"
	"os"

	"github.com/swissymissy/chipmunk/internal/auth"
)

// CLI helper that hash password and print hashespassword in terminal
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: ./hashpw <password>")
		os.Exit(1)
	}
	hash, err := auth.HashPassword(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Copy and paste this hash to PROFESSOR_PASSWORD_HASH in .env file: %s\n",hash)
}

