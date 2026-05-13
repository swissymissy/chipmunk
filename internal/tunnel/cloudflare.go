package tunnel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type QuickTunnel struct {
	Cmd       *exec.Cmd
	PublicURL string
}

type NamedTunnel struct {
	Cmd *exec.Cmd
}

// start a quick tunnel to localhost using cloudflare
func StartQuickTunnel(ctx context.Context, localURL string) (*QuickTunnel, error) {
	// find cloudflare.exe 
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exePath)
	cloudflarePath := filepath.Join(exeDir, "cloudflared.exe")
	cmd := exec.CommandContext(ctx, cloudflarePath, "tunnel", "--url", localURL)

	// catch stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	// catch stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start cloudflared: %w", err)
	}

	// channels to catch generated link, errors from stderr
	urlChan := make(chan string, 1)

	// start a scanning for tunnel url in background
	go scanForTunnelURL(stdout, urlChan)
	go scanForTunnelURL(stderr, urlChan)

	// set a timer
	timer := time.NewTimer(60 * time.Second)
	defer timer.Stop()

	select {
	case url := <-urlChan:
		return &QuickTunnel{
			Cmd:       cmd,
			PublicURL: url,
		}, nil
	case <-timer.C:
		_ = cmd.Process.Kill()
		return nil, errors.New("timeout waiting for Cloudflare tunnel URL")
	}
}

// terminate tunnel, kill the child process
func (t *QuickTunnel) Stop() error {
	if t == nil || t.Cmd == nil || t.Cmd.Process == nil {
		return nil
	}

	return t.Cmd.Process.Kill()
}


// start a named tunnel using a token from Cloudflare ZeroTrust
func StartNamedTunnel(ctx context.Context, token string) (*NamedTunnel, error) {
	// find cloudflare.exe 
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exePath)
	cloudflarePath := filepath.Join(exeDir, "cloudflared.exe")

	cmd := exec.CommandContext(ctx, cloudflarePath, "tunnel" , "run", "--token", token )
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start cloudflared: %w", err)
	}
	return &NamedTunnel{Cmd: cmd}, nil
}

// stop Named tunnel
func (t *NamedTunnel) Stop() error {
	if t == nil || t.Cmd == nil || t.Cmd.Process == nil {
		return nil
	}
	return t.Cmd.Process.Kill()
}