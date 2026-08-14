# Professor Student Management — Development Log

**Features (one session, three related changes + one ops investigation):**
1. **Student search** — a filter box over every student list so the professor can find one
   student among 40–70 without scrolling.
2. **Tab persistence** — the dashboard remembers the active tab across a page refresh.
3. **All Students directory** — a global student list so the professor can reset a password
   for a student who was removed from every course roster.
4. **The Cloudflare cache saga** — why a shipped build showed an empty tab, and how we
   diagnosed it (recurring theme; the durable fix is still a follow-up).

**Status:** All three features complete and on `main`
(`cacc3a2` search + tab persistence; `84b1adc` / `#29 3e6c4a9` All Students).
The cache-busting fix (§3.8) is **recommended but not yet implemented**.

**Working pattern this round:** search + tab persistence are frontend-only; All Students is
frontend + a thin backend handler. The backend was almost entirely pre-existing.

## 1. The Issues

### Issue 1 — Finding one student is slow
Every per-student action the professor performs lives inside a rendered list:

| Tab | List | Per-student actions |
|---|---|---|
| Roster | `#roster-list` | Reset password, Remove |
| Attendance (live) | `#attendance-list` | Mark Present / Absent |
| Edit Records (past) | `#records-list` | Mark late/present/absent, add note |

With 40–70 rows, locating a name/ID meant scrolling the whole table. Slow and error-prone.

### Issue 2 — A page refresh throws away the current tab
The dashboard is a tab-based SPA. The active tab lived only in the DOM (the `active` class on
the clicked button + the default-visible panel). A refresh always snapped back to the default
**Courses** tab — especially annoying on **Edit Records**, where the professor then had to
click back and re-select course/session.

### Issue 3 — A removed student is unreachable
"Remove" from a roster only **unenrolls** a student — the account (and its password) is kept.
But the only **Reset password** button lived inside the per-course roster. So once a student
was removed from *every* course, they no longer appeared in any list, and a professor had
**no way to reset their password** when they forgot it. This is the concrete case that
triggered the feature: a removed student forgot their password.

---

## 2. The Plan

### Issue 1 — client-side filter, no backend
The lists are already fully loaded in the browser (40–70 rows) and all three share the same
first two columns — **Student ID (col 0)** and **Name (col 1)**. So this is a *filtering*
problem, not a search-the-database problem. A server-side `LIKE` query would only pay off
with pagination or huge lists; here it would add a network round-trip and more code for no
benefit. Decision: **one shared client-side helper**, a search box above each list.

### Issue 2 — persist the tab in `localStorage`
Save the active tab on every switch, restore it on load. Chose `localStorage` (invisible,
matches how `professor_token` is already stored) over a URL hash. Scope limited to the
**tab**, not the selected course/session inside it.

### Issue 3 — a global "All Students" list
Reuse everything that already exists:
- the `ListAllStudents` query (already written and generated, previously unused),
- the reset-password endpoint (already exists, already keyed on the UUID),
- the search box and tab persistence from Issues 1 & 2.

Only two genuinely new pieces: a list handler + route, and a new tab.

---

## 3. Challenges & How We Solved Them

### 3.1 One helper for three lists — the shared-column invariant
Roster, Attendance, and Edit Records are all built by the same `buildTable`, and all put
**Student ID in column 0 and Name in column 1**. That invariant let a single function serve
every list:

```js
function filterStudentTable(containerId, query) {
    const q = query.trim().toLowerCase();
    const table = document.getElementById(containerId).querySelector("table");
    if (!table) return;
    for (const row of table.querySelectorAll("tr")) {
        const cells = row.querySelectorAll("td");
        if (cells.length === 0) continue;            // skip the header row (th only)
        const hay = (cells[0].textContent + " " + cells[1].textContent).toLowerCase();
        row.style.display = hay.includes(q) ? "" : "none";
    }
    // (also updates an optional "<containerId>-count" element → "Showing X of Y")
}
```

**Why match only cells 0–1** and not the whole row: matching `row.textContent` would hit
button labels ("Mark **Present**"), the status column, notes, and emails — e.g. typing
"present" would "match" every row. Restricting to ID + Name is what the professor actually
searches by, and it keeps results clean.

### 3.2 The re-render gotcha (the non-obvious part of search)
Filtering by toggling `row.style.display` is fine — until the table is **rebuilt**, which
resets every row to visible. Two lists rebuild out from under the filter:

- **Attendance polls every 5 seconds** (`loadAttendanceRoster`) — without care, the
  professor's filter would silently clear itself every 5s.
- **Roster / Edit Records** rebuild whenever the professor changes course or session.

Fix: a `reapplyFilter(searchId, containerId)` helper called at the **end of each render**
(`loadRoster`, `loadAttendanceRoster`, `loadRecordEditor`, and later `loadAllStudents`), so
an active query re-applies to the freshly built table.

### 3.3 Restoring a tab when there is no click to read
`showTab(tabName, btn)` used the clicked `btn` to add the `active` class. On a refresh-restore
there is no click. Solution: tag each sidebar button with `data-tab="…"`, make `btn` optional,
and look it up when missing:

```js
if (!btn) btn = document.querySelector('.tab-btn[data-tab="' + tabName + '"]');
if (btn) btn.classList.add("active");
localStorage.setItem("dashboard_tab", tabName);   // remember for next refresh
```

On load, restore the saved tab, guarding against a stale/removed name so a bad value can't
crash the page:

```js
const savedTab = localStorage.getItem("dashboard_tab");
if (savedTab && document.getElementById("tab-" + savedTab)) showTab(savedTab);
```

### 3.4 Don't leak the password hash (All Students)
`ListAllStudents` is `SELECT * FROM students`, so each returned row carries
`password_hash`, `verified`, and `must_change_password`. Returning the raw DB rows would
**leak the hashes**. The handler copies only the six safe fields into the API `Student`
model — the same pattern `HandlerRosters` already used:

```go
list := make([]Student, 0, len(students))
for _, s := range students {
    list = append(list, Student{
        ID: s.ID, StudentID: s.StudentID, Email: s.Email,
        FirstName: s.FirstName, LastName: s.LastName, Specialty: s.Specialty.String,
    })
}
```

### 3.5 The reset flow already worked for removed students — because it keys on the UUID
The existing `PUT /api/students/{student_id}/reset-password` takes the **UUID `students.id`**,
not course enrollment. So the All Students tab needed **zero new action code** — its button
calls the same `resetStudentPassword(s.id, name)` the roster tab uses. Enrollment is
irrelevant to the reset; the only thing that was ever missing was a *way to see* the student.
(This is the same `id`-vs-`student_id` data-model fact from the profile-feature log: internal
UUID for identity, school ID for humans.)

### 3.6 Route placement — no ServeMux collision
`GET /api/students` is an **exact** path. It does not conflict with `GET
/api/students/myprofile` (different path) or `PUT /api/students/{student_id}/reset-password`
(different method + a path segment). Registered under `RequireProfessor`.

### 3.7 The Cloudflare cache saga — "the tab is there but empty"
We shipped `chipmunk-v1.3.exe` (which included All Students) **without purging Cloudflare**,
deliberately, to test whether the `no-cache` header we'd added was enough. On the professor's
machine the **All Students tab appeared but listed no students**, and her terminal showed
scary red tunnel lines:

```
ERR failed to run the datagram handler error="context canceled"
WRN control stream encountered a failure while serving
INF Retrying connection in up to 2s
```

**Red herring:** those are normal Cloudflare-tunnel **QUIC connection churn** — a connection
dropped and re-registered 2s later; the log ended with a fresh `Registered tunnel connection`
and `Professor has logged in`. The tunnel was healthy. Not the cause.

**The real diagnosis, from the symptom shape:** the tab *button* appeared (so the new
`dashboard.html` reached her) but the list never populated (so the new `dashboard.js`, which
defines `loadAllStudents`, was **not** running). HTML and JS ship embedded in the **same
.exe** — if the JS were old, the HTML would be old too and the button wouldn't exist. Fresh
HTML + stale JS can therefore only mean **a cache is serving an old `dashboard.js`** — and
Cloudflare edge-caches `.js` but **not** `.html` by default. Classic.

**Isolating which cache:**
- Purge Cloudflare + open in **incognito** → **works**. So the purge fixed the *edge*, and
  fresh content is being served.
- **Normal browser mode**, even after "clear cache" → **still stale**. Same machine, same
  network, same edge — the only difference from incognito is the **local browser cache**. So
  the stale `dashboard.js` was now living in her browser, and her "clear cache" hadn't evicted
  it (commonly: served from memory cache with the tab open, or a time-range that missed it).

**Fixes applied:** purge Cloudflare; then, on her machine, "Empty Cache and Hard Reload"
(DevTools open → right-click reload) or fully quit/reopen the browser. That populated the tab.

**Diagnostics worth remembering:**
- `cf-cache-status` response header on `dashboard.js`: `HIT` = served stale from edge;
  `MISS`/`REVALIDATED`/`DYNAMIC` = went to origin / not edge-cached.
- **Incognito vs normal** cleanly separates *edge* cache from *browser* cache.
- **"Tab shows but stays empty"** is the specific tell of stale JS behind fresh HTML.
- Check the actual `Cache-Control` the browser receives — if Cloudflare rewrites `no-cache`
  to `public, max-age=…`, browsers will keep caching.

### 3.8 Why `no-cache` wasn't enough (and the durable fix — not yet done)
The static handler sends `Cache-Control: no-cache` (`middleware.NoCache` on `mux.Handle("/",
…)`). `no-cache` means "store but **revalidate** before serving," not "don't store." Two
things undermined it here:
1. **Cloudflare isn't reliably honoring it** for `.js` (a default extension it caches),
   so the edge served stale copies until purged.
2. The embedded frontend (`//go:embed`) has **no ETag/Last-Modified** — `embed.FS` files have
   a zero modtime — so there's no validator; revalidation just re-downloads a full `200`
   (harmless, but no `304` savings), and browsers that stored an older copy under permissive
   headers won't revalidate at all.

Net: **`no-cache` + manual purge is not a sustainable delivery story.** Recommended permanent
fix, in order of preference:
- **Version the asset URLs** (`/js/dashboard.js?v=1.3`), bumped per release — matches the
  existing `vX.Y.exe` versioning. A new URL means both Cloudflare and every browser fetch
  fresh automatically. No purges, no cache-clearing, ever. (Touches the `<script>`/`<link>`
  tags across the HTML files, so pair it with a one-command bump or a tiny build step.)
- **Cloudflare Cache Rule** on `/js/*` and `/css/*` → "Bypass cache" / "Respect origin
  headers," so `no-cache` actually passes through. No code, but a dashboard change.

This connects to the earlier **login-loop** incident (see `student_profile_feature.md` §4.11)
— same class of problem. Adding `no-cache` since then reduced but did **not** eliminate it.

---

## 4. What We Had Before vs. What We Added

| Area | Before | Added this session |
|---|---|---|
| Finding a student | Scroll a 40–70 row table in Roster / Attendance / Edit Records | `filterStudentTable` + `reapplyFilter`; a search box + "Showing X of Y" count above each list |
| Page refresh | Always reset to the Courses tab | `data-tab` on every tab button; `showTab` persists to `localStorage`; restore on load |
| Managing a removed student | Reachable only via a per-course roster → impossible once removed from all courses | `GET /api/students` + `HandlerListAllStudents`; an "All Students" tab with a Reset-password button (reuses the existing UUID-keyed reset) |
| Static-asset delivery | `Cache-Control: no-cache`, purge Cloudflare each release | Diagnosed that this is insufficient; documented the durable fix (still to implement) |

---

## 5. Final Architecture

### Backend (All Students only)
| Piece | File | Notes |
|---|---|---|
| Query `ListAllStudents :many` | `sql/queries/students.sql` | **Pre-existing** (`SELECT * FROM students`), was unused; now consumed |
| Handler `HandlerListAllStudents` | `internal/handlers/handler_list_all_students.go` | **New**; maps DB rows → API `Student` model (drops hash/verified/must_change_password) |
| Route `GET /api/students` | `cmd/server/main.go` | **New**; `RequireProfessor`, exact path |
| Reset endpoint | `cmd/server/main.go` | **Unchanged**; `PUT /api/students/{student_id}/reset-password`, UUID-keyed |

### Frontend
| File | Change |
|---|---|
| `js/dashboard.js` | `filterStudentTable`, `reapplyFilter`; `reapplyFilter` calls in `loadRoster`/`loadAttendanceRoster`/`loadRecordEditor`/`loadAllStudents`; `showTab` persists + restores tab; `loadAllStudents`; wire `students` tab |
| `dashboard.html` | Search input + count span above each of the 4 lists; `data-tab` on the 7 tab buttons; **All Students** sidebar button + `#tab-students` panel |
| `css/dashboard.css` | `.student-search` (capped width) and `.search-count` |
| `internal/middleware/cache.go` | Comment corrected (the header governs edge + browser, not just "client side") |

### Data-model reminder
`students.id` (UUID PK) is identity everywhere internally; `students.student_id` (school ID)
and `email` are human-facing `UNIQUE` fields. Enrollments/attendance FK to the UUID — which is
why password reset works regardless of enrollment.

---

## 6. Known Follow-ups / Deferred Work
- **Cache-busting (§3.8).** Version the JS/CSS URLs (or add a Cloudflare Cache Rule). This is
  the one that ends the purge-every-release cycle; **not yet implemented.**
- **Re-enroll a removed student.** There is currently **no professor-side "add student to a
  course"** at all — students self-enroll from their profile. The All Students tab is
  read + reset-password only; re-enrollment is a separate, still-open gap.
- **Delete account from All Students.** The `DeleteStudent` query exists but was left off the
  quick fix (destructive).
- **Restore selected course/session on refresh.** Tab persistence restores the tab but not the
  in-tab selection (e.g. which session was open in Edit Records).

---

## 7. Testing Checklist
- Search: in Roster / Attendance / Edit Records / All Students, typing filters by name **and**
  ID; the "Showing X of Y" count updates; button-label text ("Mark Present") does **not** match.
- Search survives re-render: type a query in **Attendance**, wait past the 5s poll — the filter
  stays applied. Change course in Roster/Records — filter re-applies to the new table.
- Tab persistence: switch to **Edit Records**, refresh → still on Edit Records (dropdowns reset,
  which is expected). First-ever load with nothing saved → Courses.
- All Students: open the tab → every registered student loads. Remove a student from all
  courses, then find them here and **Reset password** → succeeds (keyed on UUID).
- Security: `GET /api/students` response contains **no** `password_hash` / `verified` /
  `must_change_password`.
- Delivery: after a redeploy behind the tunnel, confirm `dashboard.js` is fresh
  (`cf-cache-status` not a stale `HIT`; incognito shows the new build). Until §3.8 lands,
  purge Cloudflare on each frontend release.