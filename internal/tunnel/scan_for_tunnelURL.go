package tunnel

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
)

func scanForTunnelURL(reader io.Reader, urlCh chan<- string, errCh chan<- error) {
	scanner := bufio.NewScanner(reader)

	re := regexp.MustCompile(`https://[a-zA-Z0-9-]+\.trycloudflare\.com`)

	for scanner.Scan() {
		line := scanner.Text()

		// print to debug
		fmt.Println("[cloudflared] ", line)

		match := re.FindString(line)
		if match != "" {
			urlCh <- match
			return
		}
	}

	if err := scanner.Err(); err != nil {
		errCh <- fmt.Errorf("read cloudflared output: %w", err)
		return
	}

	errCh <- errors.New("cloudflared stopped before URL was found")
}