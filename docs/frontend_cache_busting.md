# Frontend Cache-Busting (Versioned Assets) — Development Log

**Feature:** Version-stamp every local JS/CSS URL (`/js/dashboard.js?v=v1.4`) so each
build's assets are fetched fresh **automatically** — no Cloudflare purge, no browser
hard-reload, no dashboard visit. A new build "just works" for the professor.

**Status:** Implemented; `go build ./...` and `go test ./...` green. This closes the durable
fix that was deferred in `professor_student_management.md` §3.8.

---

## 1. The Issue

Every time we shipped a new frontend, the professor kept seeing the **old** version — the
login loop, then the empty "All Students" tab. Each release we had to manually **purge the
Cloudflare cache** and then talk her through an **"Empty Cache and Hard Reload"** in her
browser. That's unsustainable: she isn't going to open DevTools, and it's easy for us to
forget the purge and ship a "broken" build.

The root cause is two caches we don't control by default:
- **Cloudflare's edge** caches `.js`/`.css` by file extension and doesn't reliably honor our
  `no-cache` header, so it served stale assets until purged.
- **The browser** independently keeps its own copy and, once it has stored a file, may not
  re-check for a new one.

We wanted a fix where **a new build is picked up on its own**, at both layers, with zero
manual steps — because the professor's experience of "I got the new .exe but the page is
old" is exactly the failure we kept hitting.

---

## 2. The Plan

**Chosen approach: version the asset URLs.** Append `?v=<build-version>` to every local
`/js` and `/css` reference. Because caches key on the *full URL*, a new version is a new URL
the caches have never seen → they fetch it fresh. Then cache those versioned assets
**immutably** (a year, no revalidation) for speed, and keep the **HTML always fresh**
(`no-cache`) so it always hands out the current `?v=`.

The version is the **build version** from git (`git describe`), stamped into the binary at
build time and injected into the HTML at server startup.

**Alternatives considered and rejected:**
- **Cloudflare purge API on startup** — automates the dashboard step, but needs a scoped API
  token living on the professor's machine, and it only fixes the *edge* cache, not her
  *browser*. (The browser staleness is half of what bit us.) Versioning fixes both and needs
  no credential.
- **Hashed filenames** (`dashboard.4f3a1c.js`) — the gold standard, but needs a build step to
  rename files and rewrite references. Overkill for a one-professor app.
- **Placeholder token** (`?v=__VERSION__` hand-written in each page) — simple code, but you
  must remember to add the token to every new asset ref. We chose a string-scan that
  auto-covers any `/js`/`/css` reference instead (see §4).

---

## 3. How It Works (the mechanism)

This is the part to re-read if the "why" ever gets fuzzy.

### The one idea: the URL *is* the cache key
A cache (browser or Cloudflare edge) is a key→value store where the **key is the full request
URL** (query string included). `/js/dashboard.js` and `/js/dashboard.js?v=v1.4` are **two
different keys** — unrelated as far as any cache is concerned.

So versioning doesn't *bypass* or *disable* caching. It **changes the resource's identity**
so the cache has never heard of it → a cache MISS → a fresh fetch. When we ship v1.5, the
HTML points at `?v=v1.5`, a brand-new key, so both Cloudflare and every browser fetch the new
file. The old `?v=v1.4` entry just sits unused until it's evicted.

### The layering (why it's complete)
```
browser → /dashboard.html          (stable URL, no-cache → always fresh)
             │  reads it, sees:
             ▼
          /js/dashboard.js?v=v1.4   (versioned URL → immutable, cached ~forever)
          /css/style.css?v=v1.4
```
- **HTML = `no-cache`.** It's the entry point (reached by a fixed URL), so it can't be
  versioned — there's nothing "above" it to stamp its URL. Instead we keep it always-fresh,
  which is what lets it always emit the *current* `?v=`. HTML is tiny, so refetching it every
  load is cheap.
- **Versioned JS/CSS = `immutable, max-age=1yr`.** Safe *because the URL is versioned*: the
  bytes at `?v=v1.4` never change (a change means a new URL, `?v=v1.5`). This is faster than
  our old blanket `no-cache`, which forced a revalidation round-trip on every asset.

### Why `immutable` doesn't strand the user on an old version
`immutable` means "the bytes **at this exact URL** never change" — which is true. After an
update the browser is never asked for the old URL again (the fresh HTML points at the new
one), so it can't serve stale content for the page. `immutable` is only dangerous on an
**un-versioned** URL — which is why the `Immutable` middleware only applies it when `?v=` is
present.

### Where the version actually lives
The `?v=` rides in the **URL** (the reference), not inside the asset file. The `.js`/`.css`
files are served byte-for-byte unchanged; the version is written into the **HTML's**
`<script>`/`<link>` tags, because the HTML is the only place asset URLs are authored. That's
why the startup step only rewrites HTML.

---

## 4. What Changed

### `Makefile` — stamp the version at build time
```make
VERSION    ?= $(shell git describe --tags --always --dirty)
LDFLAGS    := -X 'main.Version=$(VERSION)' -X 'main.CommitSHA=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'
build:
	mkdir -p dist
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/chipmunk.exe ./cmd/server
```
`-ldflags -X` writes the git version into the `main.Version` variable at compile time — so the
binary knows its own version with no config file. `GOOS=windows GOARCH=amd64` cross-compiles
the Windows `.exe` from WSL.

### `cmd/server/main.go` — the version vars + the wiring
```go
var (
	Version   = "development" // overridden by -ldflags at build; fallback for `go run`
	CommitSHA = "unknown"
	BuildTime = "unknown"
)
```
```go
fileServer := http.FileServer(http.FS(frontendSub))
htmlHandler, err := HTMLHandler(frontendSub, Version, fileServer)
// ...
mux.Handle("/js/",  middleware.Immutable(fileServer)) // versioned assets: cache hard
mux.Handle("/css/", middleware.Immutable(fileServer))
mux.Handle("/",     middleware.NoCache(htmlHandler))  // HTML + images/favicon: always fresh
```
`ServeMux` routes by longest-matching prefix, so `/js/…` and `/css/…` hit `Immutable`, and
everything else (HTML pages, images, favicon) hits `NoCache(htmlHandler)`.

### `cmd/server/helpers.go` — `AppendVersion` (the injector)
Plain string scan (no regex, no placeholder) that appends `?v=<version>` before the closing
quote of every value opening with a given marker:
```go
AppendVersion(html, `"/js/`, version)   // src="/js/x.js"  -> src="/js/x.js?v=..."
AppendVersion(html, `"/css/`, version)  // href="/css/y.css" -> href="/css/y.css?v=..."
```
**Why the marker includes the leading quote** (`"/js/`): it matches only local
`attr="/js/…"` values and skips the external CDN script (`src="https://…"`) and images
(`href="/images/…"`). That selectivity is the whole trick — a `.js`-ending CDN URL like
`qrcode.min.js` is left alone because its opening quote is followed by `https`, not `/js/`.

### `cmd/server/frontend_handler.go` — `HTMLHandler` (build once, serve many)
At **startup** it walks the embedded frontend, runs each `.html` through `AppendVersion` for
`/js/` and `/css/`, and caches the stamped bytes in a `map[string][]byte`. Per **request** it
serves the pre-stamped HTML for known pages (via `http.ServeContent`) and falls back to the
file server for everything else (images, favicon). Building once means the scan runs a single
time, and the map is read-only afterward so it's concurrency-safe without locks.

### `internal/middleware/cache.go` — `NoCache` (kept) + `Immutable` (new)
```go
// Immutable: cache versioned assets for a year with no revalidation; anything
// without ?v= falls back to no-cache so an un-versioned URL is never pinned.
func Immutable(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("v") != "" {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		next.ServeHTTP(w, r)
	})
}
```
Cache-header policy lives in middleware; content lives in the handler. Clean separation, and
`NoCache` stays in use.

### `cmd/server/helpers_test.go` — table-driven tests for `AppendVersion`
11 cases: core behavior, selectivity (CDN/images/no-match/mixed left correct), graceful
failure (malformed / empty), and two characterization tests pinning current behavior (empty
version still appends `?v=`; an existing `?query` yields a double `?`).

---

## 5. Guide for Future Use

### Cutting a release
```bash
git commit -am "…"     # commit first, or the version gets a -dirty suffix
git tag v1.5           # tag the release
make build             # stamps ?v=v1.5, outputs dist/chipmunk.exe
```
Then ship `dist/chipmunk.exe`. **No purge, no cache-clearing** — the professor's browser and
Cloudflare fetch the new `?v=v1.5` assets on their own because the URLs changed.

To build without tagging (quick one-off), override the version:
```bash
make build VERSION=v1.5-test
```

### Verifying it worked
You can't run the Windows `.exe` on WSL, so test the mechanism with a native build:
```bash
go build -ldflags="-X main.Version=v1.5" -o /tmp/chipmunk ./cmd/server
# run it, then:
curl -s  localhost:8080/dashboard.html | grep '/js/'      # -> ...dashboard.js?v=v1.5
curl -sI localhost:8080/dashboard.html                    # -> Cache-Control: no-cache
curl -sI "localhost:8080/js/dashboard.js?v=v1.5"          # -> Cache-Control: public, max-age=31536000, immutable
```
Locally (via `go run`, no ldflags) the version is `development`, so you'll see
`?v=development` — that's expected.

### Gotchas / things future-me should remember
- **Commit before building**, or `git describe --dirty` appends `-dirty` to the version.
- **Building from WSL needs `GOOS=windows`** (already in the Makefile) — otherwise you get a
  Linux binary named `.exe` that won't run for the professor.
- **New asset *types* need wiring.** `AppendVersion` only stamps `/js/` and `/css/`. If you
  ever add, say, `/fonts/`, you must (a) add an `AppendVersion(…, "/fonts/", …)` call in
  `HTMLHandler`, and (b) add `mux.Handle("/fonts/", middleware.Immutable(fileServer))`.
  Individual new files under existing `/js/` or `/css/` are covered automatically.
- **Images are intentionally not versioned** — they're served `no-cache` (always revalidate).
  A logo/favicon rarely changes; if you ever need one cached hard, version it like the JS/CSS.
- **Don't give a local asset ref a `?query` in the HTML** — `AppendVersion` would produce a
  double `?` (see the characterization test).

---

## 6. Testing
`go test ./cmd/server/` runs `TestAppendVersion` (11 cases). It's picked up by the CI
`tests` job (`go test -race -cover ./...`). If you change `AppendVersion`'s behavior on the
two characterization cases (empty version, existing query), update the code **and** those
`want` values in the same commit — that's the signal it changed on purpose.

---

## 7. Known Limitations / Follow-ups
- **Double `?`** if a local ref already has a query string (documented, not hit today).
- **Images/favicon un-versioned** (deliberate; `no-cache`).
- **One global version busts everything** each build (git SHA/tag changes → all `?v=` change),
  even files that didn't change. Fine at this scale; content-hash-per-file would be the
  finer-grained upgrade if it ever matters.
- The `Immutable`/`NoCache` split assumes assets live under `/js/` and `/css/`. If the
  frontend layout changes, revisit the route prefixes in `main.go`.