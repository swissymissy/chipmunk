package main

import (
	"bytes"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
)

// the handler that reads each HTML at startup
// then stamp the version and serve
func HTMLHandler(fsys fs.FS, version string, fallback http.Handler) (http.Handler, error) {
	// build a lookup table to store stamped pages
	// key = request path, value = version-stamped
	pages := make(map[string][]byte)

	// walk through the embedded file system and process on each HTML file
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(p, ".html") {
			return nil // skip directories and non-html files
		}

		// read file content
		rawByte, err := fs.ReadFile(fsys, p) // #nosec G304
		if err != nil {
			return err
		}

		key := "/" + p
		// stamp version to asset files
		stamped := AppendVersion(string(rawByte), `"/js/`, version)
		stamped = AppendVersion(stamped, `"/css/`, version)
		pages[key] = []byte(stamped)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		if requestPath == "/" {
			requestPath = "/index.html"
		}
		if html, ok := pages[requestPath]; ok {
			http.ServeContent(w, r, path.Base(requestPath), time.Time{}, bytes.NewReader(html))
			return
		}
		fallback.ServeHTTP(w, r)
	}), nil
}
