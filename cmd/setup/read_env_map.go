package main

import (
	"os"
	"strings"
)

// reads a .env file and save values into a map
func ReadEnvMap(path string) (map[string]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	env := make(map[string]string)
	lines := strings.Split(string(content), "\n")

	for _, line := range lines{
		line := strings.TrimSpace(line)

		// skip empty lines and comments
		if line == "" || strings.HasPrefix(line,"#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// remove double/single quote from value
		value = strings.Trim(value, `"'`)
		env[key] = value
	}
	return env, nil
}

