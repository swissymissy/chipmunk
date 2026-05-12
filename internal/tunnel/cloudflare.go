package tunnel

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"time"
)


type QuickTunnel struct {
	Cmd *exec.Cmd
	PublicURL string
}

// start a quick tunnel to localhost using cloudflare
func StartQuickTunnel(ctx context.Context, localURL string) (*QuickTunnel, error) {
	cloudflarePath := "cloudflared"
	if runtime.GOOS == "windows" {
		cloudflarePath = "cloudflared.exe"
	}

	// run cloudflared tunnel --url as child process
	cmd := exec.CommandContext(ctx, cloudflarePath, "tunnel", "--url", localURL)

	// catch stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	// catch stderr
	stderr , err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start cloudflared: %w", err)
	}

	// channels to catch generated link, errors from stderr
	urlChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// start a scanning for tunnel url in background
	go scanForTunnelURL(io.MultiReader(stdout, stderr), urlChan, errChan)

	// set a timer
	timer := time.NewTimer(20 * time.Second)
	defer timer.Stop()

	select {
	case url := <-urlChan:
		return &QuickTunnel{
			Cmd: cmd,
			PublicURL: url,
		}, nil
	case err := <-errChan:
		_ = cmd.Process.Kill()
		return nil, err 
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