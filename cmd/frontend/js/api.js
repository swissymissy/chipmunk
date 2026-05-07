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
        localStorage.removeItem("professor_token");
        authToken = null;
        window.location.href = "/prof_login.html";
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

// error handler function
let errorHandler = (msg) => console.error(msg);

function setErrorHandler(fn) { errorHandler = fn; }

async function safe(fn) {
    try { await fn(); } catch (err) { errorHandler(err.message); }
}