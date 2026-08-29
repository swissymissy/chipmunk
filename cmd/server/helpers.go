package main

import (
	"strings"
)

// helper function for parsing HTML
// it will appends the string "?v=<version>" to every attribute value that opens with <open>
// for example: "/js/change_password.js" opens with '"/js/'
// this function will parse the HTML file, find any path matches the <open>, it will append the version before the closing quote
// for example, it finds the path "/js/change_password.js" opens with '"/js/', it then appends "?v=<version>" to that path
// result: "/js/change_password.js?v=<version>"
// this only works for the local file, any plugin will be left alone

func AppendVersion(html, open, version string) string {
	var b strings.Builder
	v := "?v=" + version
	// parse html
	for {
		i := strings.Index(html, open) // find the <open> substring in html

		// can't find the <open> string
		if i == -1 {
			b.WriteString(html) // write the rest then stop the loop
			return b.String()
		}

		valueStart := i + 1                            // index of first char after the opening quote
		q := strings.IndexByte(html[valueStart:], '"') // find position of the ending quote
		if q == -1 {
			b.WriteString(html) // can't find the ending quote - malformed
			return b.String()
		}

		closeQuote := valueStart + q     // index of the ending quote
		b.WriteString(html[:closeQuote]) // start appending everything up to the closing quote
		b.WriteString(v)                 // append the version
		html = html[closeQuote:]         // continure parsing from the closing quote
	}
}
