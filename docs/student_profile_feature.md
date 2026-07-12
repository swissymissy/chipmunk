# Student Profile Management — Development Log

**Feature:** Let students log in on their own and edit their profile (school ID, email,
name) and manage their course enrollments, without having to scan a check-in QR code.

**Status:** Backend complete; frontend complete. One related item (professor removing a
student from a roster) was scoped but deliberately deferred.

---

## 1. Background

Chipmunk is a GPS-based classroom attendance app. A professor starts a session on their
laptop (which captures the classroom's coordinates), students scan a QR code to reach a
check-in page, log in, and the server uses the Haversine formula to confirm the student is
physically within the classroom radius before marking them present.

After the professor used it in a real class (~45–70 students), two problems surfaced.

---

## 2. The Issues

### Issue 1 — Students can't fix their own information
If a student typed the wrong school ID or misspelled their email during registration, they
were **stuck** with it. There was no way to edit their own account. The only time a student
ever authenticated was transiently inside the check-in page (scan QR → log in → check in),
and that token was never persisted — so there was no "logged-in student" area at all.

### Issue 2 — Professor can't remove a student from a course roster
If a student's record was hopelessly messed up, the professor had no way to remove that one
student from a course so the student could re-register cleanly. (This issue was analyzed and
planned but **parked** to focus on Issue 1 first.)

---

## 3. The Plan

We deliberately split the work: **backend first** (written and reviewed field-by-field),
then **frontend**.

### Scope decided for this round
- Students can log in **without scanning** a QR code.
- Students can edit **school ID**, **email**, and **name** (name was added mid-way).
- Students can **add / remove courses**.
- Password editing and Issue 2 were pushed to a later round.

### Backend design
- New per-field SQL queries (see §5) keyed on the student's UUID.
- New handlers, one endpoint per editable field.
- Standard guardrails: validate input, key every mutation on the JWT identity, don't leak
  the password hash.

### Frontend design
- A **generic landing page** at `/` with two buttons — *Professor login* and *Student
  login* — so the two roles never collide.
- A dedicated **student login page** and **profile page**.
- Reuse the professor's existing login/token pattern for the student side.

---

## 4. Challenges & How We Solved Them

This is the heart of the log — the non-obvious problems and the reasoning behind each fix.

### 4.1 "student_id" means two different things
The single most important discovery. The `students` table has **two** identifiers:

| Column | Meaning | Constraints |
|---|---|---|
| `students.id` | Internal UUID, the real primary key | `PRIMARY KEY` |
| `students.student_id` | The human **school ID** (e.g. `U12345678`) | `NOT NULL UNIQUE` |

Crucially, **`enrollments.student_id` and `attendance_records.student_id` are foreign keys
to `students.id` (the UUID), *not* to the school ID** — despite the confusing column name.

**Why this mattered:** it meant letting a student change their school ID is *safe*. Nothing
references the school ID except its own `UNIQUE` constraint, so there are no cascading
updates or orphaned enrollment/attendance rows to worry about. This removed the scariest
part of the feature.

### 4.2 Always key mutations on the JWT identity, never on client input
The pre-existing `UpdateStudentEmail`/`UpdatePassword` queries keyed off `WHERE student_id
= ?` (the school ID). Once the school ID itself becomes editable, keying updates on it is
fragile. So every new profile query keys on the **UUID `id`**, which we read from the JWT
(`middleware.GetUserID`) — never from a value the client can supply. This also prevents
IDOR (a student trying to edit someone else's row).

### 4.3 Per-field endpoints vs. one combined update
We considered three shapes:
1. One `UpdateStudentProfile` that always writes every field.
2. One partial-update query using `COALESCE(sqlc.narg(...), col)`.
3. **Separate query + endpoint per field.** ← chosen

We chose per-field endpoints because they're **simpler to read** and make **409 conflict
handling trivial** — each endpoint touches exactly one `UNIQUE` column, so if a duplicate
occurs, we already know which field collided without parsing the DB error. (The write volume
is identical either way, so throughput was never the deciding factor.)

### 4.4 UNIQUE conflicts should be 409, not 500
`student_id` and `email` are `UNIQUE`. If a student picks one that already belongs to
someone else, the DB raises a constraint error — which is *not* `sql.ErrNoRows`, so a naive
handler falls through to a generic 500. We detect it and return **409 Conflict** with a
clear message ("email is already in use"). Implemented for email; **deliberately deferred
for school ID** to keep the first pass simple (a duplicate school ID currently surfaces as a
generic 500 — a known, accepted trade-off).

### 4.5 Wrong status codes caused a hidden logout bug
An early version returned **401 Unauthorized** for generic DB errors. That's wrong on its
own (a DB hiccup isn't an auth failure), but it was *also* dangerous: the frontend's
`api.js` treats **any 401 as "session dead"** — it clears the token and redirects to login.
So a transient DB error mid-edit would have **bounced the student out of their session**.
Fixed by returning **500** for real server errors and reserving 401 for genuine auth
problems.

### 4.6 A double-response bug from a missing `return`
The profile handler wrote a 404 on `sql.ErrNoRows` but was missing a `return`, so it fell
through and wrote a second response (401). That produces a `superfluous WriteHeader` warning
and a garbled body. Fixed by returning immediately after the 404.

### 4.7 `index.html` was secretly the professor dashboard
When we went to build the "generic home page," we found that `index.html` **was not a
landing page — it was the entire professor dashboard** (it loads `dashboard.js`, which
guards on the professor token). The server serves it at `/`.

Naively making `/` a landing page would have broken the professor flow. Solution:
- `git mv index.html → dashboard.html` (history preserved).
- Create a brand-new `index.html` as the actual landing page.
- Update `prof_login.js` to redirect to `/dashboard.html` after login.

### 4.8 `api.js`'s 401 handling was hard-wired to the professor
The shared `api.js` 401 branch called `clearProfessorAuth()` → removed `professor_token` →
redirected to `/prof_login.html`. If a student page reused it, a student 401 would wrongly
nuke the *professor's* session.

Solution: generalize it, mirroring the existing `setErrorHandler` pattern.
- Added `setUnauthorizedHandler(fn)` + a default no-op `onUnauthorized`.
- Each authed page registers its own: professor → `clearProfessorAuth`; student →
  `clearStudentAuth` (clears a separate `student_token` key → `/student_login.html`).
- The two sessions now use **separate localStorage keys** and never clobber each other.
- The QR check-in page is unaffected: it holds its token in a local variable and never sets
  `authToken`, so the `&& authToken` guard means its 401 path never fires.

### 4.9 Token lifetime was too short for editing
The login JWT expired in **15 minutes** — fine for a quick scan-and-check-in, but a student
could get logged out mid-edit on the profile page. We bumped `MakeJWT` to **30 minutes** (a
middle ground; an hour felt too long). The same `POST /api/auth/login` issues both
check-in and profile tokens, so it's a one-line change.

### 4.10 Fields were editable by accident
The first profile page left every field open for editing, so a student could click into a
field and change it without meaning to. We reworked each field to be **read-only by
default**: a single **Edit** button unlocks the input(s) and reveals **Save** / **Cancel**.
Save persists (and re-locks on success; stays open on error so they can fix it); Cancel
reverts to the original value with no request. One reusable controller
(`setupEditableField`) drives all three field groups, including the two-input name group.

---

## 5. Final Architecture

### Data model (unchanged, but now understood)
- `students.id` (UUID PK) is the identity used everywhere internally.
- `students.student_id` (school ID) and `students.email` are `UNIQUE` and freely editable.
- `enrollments` / `attendance_records` reference `students.id`.

### New SQL queries
| Query | File | Purpose |
|---|---|---|
| `GetProfileByID` | `students.sql` | Load a student's profile by UUID |
| `UpdateStudentSchoolID` | `students.sql` | Update school ID (keyed on UUID) |
| `UpdateStudentEmailByID` | `students.sql` | Update email (keyed on UUID) |
| `UpdateStudentName` | `students.sql` | Update first/last name |
| `RemoveACourse` | `enrollments.sql` | Delete one enrollment row |

### New / changed endpoints (all `AuthRequired`)
| Method | Route | Handler |
|---|---|---|
| `GET` | `/api/students/myprofile` | `HandlerGetStudentProfile` |
| `PUT` | `/api/students/myprofile/schoolid` | `HandlerStudentUpdateSchoolID` |
| `PUT` | `/api/students/myprofile/email` | `HandlerStudentUpdateEmail` |
| `PUT` | `/api/students/myprofile/name` | `HandlerStudentUpdateName` |
| `DELETE` | `/api/students/myprofile/courses/{id}` | `HandlerStudentRemoveACourse` |

Handlers live in `handler_students_profile.go` (read) and
`handler_students_edit_profile.go` (the four edits). Input is validated by `EmailCheck`,
`NameCheck`, and `SchoolIDCheck` in `helpers.go`.

### Frontend files
| File | Change |
|---|---|
| `index.html` | **New** landing page (Professor / Student login) |
| `dashboard.html` | Renamed from the old `index.html` (professor dashboard) |
| `js/prof_login.js` | Redirects to `/dashboard.html` |
| `js/api.js` | Generalized 401 handling; `clearStudentAuth` + `setUnauthorizedHandler` |
| `js/dashboard.js` | Registers `clearProfessorAuth` as its 401 handler |
| `student_login.html` + `js/student_login.js` | **New** student login (stores `student_token`) |
| `profile.html` + `js/profile.js` | **New** profile page (Edit/Save/Cancel + course add/remove) |
| `js/register.js` | Shows success, then redirects to `/student_login.html` |
| `checkin.html` | Added a "Manage my profile" link |
| `css/style.css` | Landing, profile-field, and edit-action styles |
| `internal/auth/jwt.go` | Token TTL 15 min → 30 min |

### Auth flow after the changes
```
                         /  (landing)
                        /            \
             Professor login       Student login
                  |                      |
           /prof_login.html      /student_login.html
                  |                      |
          professor_token         student_token
                  |                      |
           /dashboard.html         /profile.html
```

---

## 6. Known Follow-ups / Deferred Work

- **Issue 2 (professor removes a student from a roster).** Analyzed and planned; not built.
  It reuses the same `RemoveACourse`/enrollment-delete query behind a professor-only route.
- **School ID 409.** Duplicate school ID currently returns a generic 500 instead of a clean
  409 — deliberately deferred; easy to add the same UNIQUE-detection branch used for email.
- **Password change.** Out of scope this round; slots in as `PUT /api/students/myprofile/password`
  plus a form section.
- **`SchoolIDCheck` polish.** The leading `U` isn't actually enforced (only that the
  remaining chars are digits), and the length error message says "too long" even when the ID
  is too short.

---

## 7. Testing Checklist

- Professor: login → lands on `/dashboard.html`; existing flows unchanged.
- Student: register → success message → redirected to `/student_login.html`.
- Student: login → `/profile.html` loads with prefilled, **read-only** fields.
- Edit each field: Edit → Save (persists, re-locks), and Edit → Cancel (reverts).
- Duplicate email edit → inline **409** "email is already in use".
- Add a course, remove a course → list updates.
- Let the 30-min token expire mid-session → student is bounced to `/student_login.html`
  (and the professor session is untouched).
