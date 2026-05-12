package tunnel

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

func scanForTunnelURL(reader io.Reader, urlCh chan<- string) {
	scanner := bufio.NewScanner(reader)

	re := regexp.MustCompile(`https://[a-zA-Z0-9-]+\.trycloudflare\.com`)

	for scanner.Scan() {
		line := scanner.Text()

		// print to debug
		fmt.Println("[cloudflared] ", line)

		match := re.FindString(line)
		if match != "" {
			select {
			case urlCh <- match:
			default:
			}
			return
		}
	}
}
