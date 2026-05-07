let authToken = null;
function setAuthToken(t) { authToken = t; }

// api calling helper
async function api(method, url, body, token) {
    const opts = { method, headers: { "Content-Type": "application/json" } };
    if (body) opts.body = JSON.stringify(body);
    const tk = token || authToken;
    if (tk) opts.headers["Authorization"] = "Bearer " + tk;

    const res = await fetch(url, opts);

    // dashboard auth expired or missing -> bounce to login
    if (res.status === 401 && authToken) {
        clearProfessorAuth();
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

// clear stored professor auth and redirect to login.
// used by 401 handlers and the logout button.
function clearProfessorAuth() {
    localStorage.removeItem("professor_token");
    authToken = null;
    window.location.href = "/prof_login.html";
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
