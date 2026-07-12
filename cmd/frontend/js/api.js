let authToken = null;
function setAuthToken(t) { authToken = t; }

// api calling helper
async function api(method, url, body, token) {
    const opts = { method, headers: { "Content-Type": "application/json" } };
    if (body) opts.body = JSON.stringify(body);
    const tk = token || authToken;
    if (tk) opts.headers["Authorization"] = "Bearer " + tk;

    const res = await fetch(url, opts);

    // auth expired or missing -> let the current page decide where to bounce.
    // each authed page registers its own handler via setUnauthorizedHandler
    // (professor -> prof login, student -> student login) so a student 401
    // never clears the professor session and vice versa.
    if (res.status === 401 && authToken) {
        onUnauthorized();
        return;
    }
    if (!res.ok) {
        const text = await res.text();
        let msg = "Request failed";
        try { msg = JSON.parse(text).error || msg; } catch {}
        throw new Error(msg);
    }
    if (res.headers.get("Content-Type")?.includes("json")) return res.json();
}

// on-unauthorized handler: what a 401 does depends on which page we're on.
// defaults to a no-op; authed pages register their own via setUnauthorizedHandler.
let onUnauthorized = () => {};
function setUnauthorizedHandler(fn) { onUnauthorized = fn; }

// clear stored professor auth and redirect to login.
// used by the professor 401 handler and the logout button.
function clearProfessorAuth() {
    localStorage.removeItem("professor_token");
    authToken = null;
    window.location.href = "/prof_login.html";
}

// clear stored student auth and redirect to student login.
// used by the student 401 handler and the profile logout button.
function clearStudentAuth() {
    localStorage.removeItem("student_token");
    authToken = null;
    window.location.href = "/student_login.html";
}

// shared course label format used across dashboard, register, checkin.
function courseLabel(c) {
    return c.course_name + " — " + c.section_date + " " + c.start_time;
}

// populate a <select> with the given items. preserves prior selection
// if it's still in the new list.
function fillDropdown(selectId, items, valueFn, labelFn, placeholder = "-- Select --") {
    const select = document.getElementById(selectId);
    const current = select.value;
    select.innerHTML = `<option value="">${placeholder}</option>`;
    for (const item of items) {
        const opt = document.createElement("option");
        opt.value = valueFn(item);
        opt.textContent = labelFn(item);
        select.appendChild(opt);
    }
    if (current) select.value = current;
}

// disable the form's submit button while fn runs; route errors to errorHandler.
// use this from form submit listeners instead of bare safe(...) so users
// can't double-click and create duplicates.
async function submitForm(form, fn) {
    const btn = form.querySelector("button[type=submit]");
    if (btn) btn.disabled = true;
    try { await fn(); }
    catch (err) { errorHandler(err.message); }
    finally { if (btn) btn.disabled = false; }
}

// error handler function
let errorHandler = (msg) => console.error(msg);

function setErrorHandler(fn) { errorHandler = fn; }

async function safe(fn) {
    try { await fn(); } catch (err) { errorHandler(err.message); }
}


// returns a device fingerprint hash via ThumbmarkJS, or "" if the library
// isn't loaded or fingerprinting fails. callers treat "" as "no fingerprint" —
// the server stores NULL and skips flag detection for that row.
async function getDeviceFingerprint() {
    try {
        if (typeof ThumbmarkJS === "undefined") return "";
        const tm = new ThumbmarkJS.Thumbmark();
        const result = await tm.get();
        return result?.thumbmark || "";
    } catch {
        return "";
    }
}