package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
	"github.com/swissymissy/chipmunk/internal/auth"
)

// CLI helper that hash password and print hashespassword in terminal
func main() {
	runSetup()
}

//
func runSetup() {
	// ensure .env exist -> copy frm .env.example if missing
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		if err := copyFile(".env.example", ".env"); err != nil {
			fmt.Fprintf(os.Stderr, "no .env and no .env.exmaple found: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Create .env frm .env.example")
	}

	// generate JWT_SECRET if empty
	env, err := ReadEnvMap(".env")
	if env["JWT_SECRET"] == "" || env["JWT_SECRET"] == `""` {
		secret, err := generateSecret()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate jwt secret: %v\n", err)
			os.Exit(1)
		}

		if err := patchEnvLine(".env", "JWT_SECRET" , `"`+secret+`"`); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save jwt secret: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Generated JWT_SECRET")
	}

	// prompt for password
	pw := promptPasswordTwice()

	// hash pw
	hash, err := auth.HashPassword(pw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to hash password: %v\n", err)
		os.Exit(1)
	}

	// patch .env file with single-quote hash (double quote causes error)
	if err := patchEnvLine(".env", "PROFESSOR_PASSWORD_HASH" , "`"+hash+"`"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to save professor password hash: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Professor's password set")
	fmt.Println("Setup is completed. You can now run chipmunk.exe")
}

// prompt for entering password twice
func promptPasswordTwice() string {
	fmt.Print("Enter new password: ")
	pw1, err := term.ReadPassword(int(os.Stdin.Fd())) // read input from term
	if err != nil {
		fmt.Println("error reading password")
		return ""
	}
	fmt.Println("Confirm password: ")
	pw2 , err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("error reading confirmation password")
		return ""
	} 
	
	// compare 2 passwords
	if string(pw1) != string(pw2) {
		fmt.Println("password 1 and password 2 do not match")
		return ""
	}

	return string(pw1)
}

// replace old value with new value of the given key in .env file
func patchEnvLine(path, key, value string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// split into lines
	lines := strings.Split(string(content), "\n")
	found := false 
	newLine := key + "=" + value

	// find the key inside the .env file to change the value
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, key+"=") {
			lines[i] = newLine
			found = true 
			break
		}
	}
	if !found {
		lines = append(lines, newLine)
	}
	result := strings.Join(lines, "\n")

	// atomic write: write to temp, then rename
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(result), 0600); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}

// generate jwt secret
func generateSecret() string {
	b := make([]byte, 64)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

