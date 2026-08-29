package main

import "testing"

func TestAppendVersion(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		open    string
		version string
		want    string
	}{
		// --- core behavior ---
		{
			name:    "single match",
			html:    `<script src="/js/app.js"></script>`,
			open:    `"/js/`,
			version: "v1",
			want:    `<script src="/js/app.js?v=v1"></script>`,
		},
		{
			name:    "multiple matches are all stamped",
			html:    `<script src="/js/a.js"></script><script src="/js/b.js"></script>`,
			open:    `"/js/`,
			version: "v1",
			want:    `<script src="/js/a.js?v=v1"></script><script src="/js/b.js?v=v1"></script>`,
		},
		{
			name:    "works for css via href",
			html:    `<link rel="stylesheet" href="/css/style.css">`,
			open:    `"/css/`,
			version: "v2",
			want:    `<link rel="stylesheet" href="/css/style.css?v=v2">`,
		},

		// --- selectivity: only the target local prefix is touched ---
		{
			name:    "external CDN url is left alone",
			html:    `<script src="https://cdn.jsdelivr.net/npm/qrcode.min.js"></script>`,
			open:    `"/js/`,
			version: "v1",
			want:    `<script src="https://cdn.jsdelivr.net/npm/qrcode.min.js"></script>`,
		},
		{
			name:    "image path is left alone",
			html:    `<img src="/images/logo.png">`,
			open:    `"/js/`,
			version: "v1",
			want:    `<img src="/images/logo.png">`,
		},
		{
			name:    "no match returns input unchanged",
			html:    `<p>no assets here</p>`,
			open:    `"/js/`,
			version: "v1",
			want:    `<p>no assets here</p>`,
		},
		{
			name:    "mixed page stamps only the target prefix",
			html:    `<link href="/css/s.css"><script src="/js/a.js"></script><img src="/images/l.png">`,
			open:    `"/js/`,
			version: "v1",
			want:    `<link href="/css/s.css"><script src="/js/a.js?v=v1"></script><img src="/images/l.png">`,
		},

		// --- edge cases / graceful failure ---
		{
			name:    "malformed: missing closing quote returns input unchanged",
			html:    `<script src="/js/a.js`,
			open:    `"/js/`,
			version: "v1",
			want:    `<script src="/js/a.js`,
		},
		{
			name:    "empty input",
			html:    "",
			open:    `"/js/`,
			version: "v1",
			want:    "",
		},

		// --- characterization: these pin current behavior; revisit if undesired ---
		{
			name:    "empty version still appends ?v=",
			html:    `<script src="/js/a.js"></script>`,
			open:    `"/js/`,
			version: "",
			want:    `<script src="/js/a.js?v="></script>`,
		},
		{
			name:    "existing query yields a double question mark (known limitation)",
			html:    `<script src="/js/a.js?foo=1"></script>`,
			open:    `"/js/`,
			version: "v1",
			want:    `<script src="/js/a.js?foo=1?v=v1"></script>`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := AppendVersion(tc.html, tc.open, tc.version)
			if got != tc.want {
				t.Errorf("AppendVersion(%q, %q, %q)\n got: %q\nwant: %q",
					tc.html, tc.open, tc.version, got, tc.want)
			}
		})
	}
}
