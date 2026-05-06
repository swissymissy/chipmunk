// api calling helper
async function api(method, url, body, token) {
    const opts = { method, headers: { "Content-Type": "application/json" } };
    if (body) opts.body = JSON.stringify(body);
    if (token) opts.headers["Authorization"] = "Bearer " + token;
    const res = await fetch(url, opts);
    if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || "Request failed");
    }
    if (res.headers.get("Content-Type")?.includes("json")) return res.json();
}

// error handler function
let errorHandler = (msg) => console.error(msg);

function setErrorHandler(fn) { errorHandler = fn; }

async function safe(fn) {
    try { await fn(); } catch (err) { errorHandler(err.message); }
}