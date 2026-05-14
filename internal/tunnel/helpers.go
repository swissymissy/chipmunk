package tunnel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// helper that tries windows-location first, the falls back to PATH for linux
// it returns the path to the cloudflared binary
// windows: look next to the running binary (chipmink.exe + cloudflared.exe in same folder)
// dev: fallback to PATH
func resolveCloudflaredPath() (string, error) {
	binName := "cloudflared"
	if runtime.GOOS == "windows" {
		binName = "cloudflared.exe"
	}

	// using absolute path to find the exectable file in windows
	if exePath, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(exePath), binName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	// linux (PATH)
	if found, err := exec.LookPath(binName); err == nil {
		return found, nil
	}

	return "", fmt.Errorf("cloudflared not found (looked next to binary and on PATH)")
}
