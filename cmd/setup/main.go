package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/swissymissy/chipmunk/internal/auth"
)

// setup file to let professor set up jwt secret and password
func main() {
	runSetup()
}

func runSetup() {
	// ensure .env exist -> copy frm .env.example if missing
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		if err := CopyFile(".env.example", ".env"); err != nil {
			fmt.Fprintf(os.Stderr, "no .env and no .env.exmaple found: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Create .env frm .env.example")
	}

	// generate JWT_SECRET if empty
	env, err := ReadEnvMap(".env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read .env: %v\n", err)
		os.Exit(1)
	}
	if env["JWT_SECRET"] == "" || env["JWT_SECRET"] == `""` {
		secret, err := generateSecret()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate jwt secret: %v\n", err)
			os.Exit(1)
		}

		if err := PatchEnvLine(".env", "JWT_SECRET" , `"`+secret+`"`); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save jwt secret: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Generated JWT_SECRET")
	}

	// prompt for password
	pw := PromptPasswordTwice()

	// hash pw
	hash, err := auth.HashPassword(pw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to hash password: %v\n", err)
		os.Exit(1)
	}

	// patch .env file with single-quote hash (double quote causes error)
	if err := PatchEnvLine(".env", "PROFESSOR_PASSWORD_HASH" , "'"+hash+"'"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to save professor password hash: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Professor's password set")
	fmt.Println("Setup is completed. You can now run chipmunk.exe")
}



// generate jwt secret
func generateSecret() (string, error) {
	b := make([]byte, 64)
	
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

