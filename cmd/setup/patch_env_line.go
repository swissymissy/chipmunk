package main

import (
	"os"
	"strings"
)

// replaces or appends KEY=value in .env file
func PatchEnvLine(path, key, value string) error {
	content, err := os.ReadFile(path) // #nosec 304 -- local setup tool, operator-provided path
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	found := false
	newLine := key + "=" + value

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// skip comments
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

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

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(result), 0600); err != nil { // #nosec G703 -- local setup tool, operator-provided path
		return err
	}

	return os.Rename(tmp, path)
}
