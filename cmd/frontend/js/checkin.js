let jwtToken = null;
let qrToken = null;

document.addEventListener("DOMContentLoaded", () => {
    qrToken = new URLSearchParams(window.location.search).get("t");
    if (!qrToken) { showCheckinError("No check-in code found. Please scan the QR code."); return; }
    document.getElementById("login-form").addEventListener("submit", e => {
        e.preventDefault();
        safe(handleLogin);
    });
});

async function handleLogin() {
    clearError();
    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value;

    const data = await api("POST", "/api/auth/login", { email, password });
    jwtToken = data.token;

    document.getElementById("greeting").textContent = "Hi " + data.first_name + "!";
    showSection("loading-section");
    startCheckIn();
}

function startCheckIn() {
    if (!navigator.geolocation) {
        showCheckinError("Your browser does not support location services.");
        return;
    }
    navigator.geolocation.getCurrentPosition(
        pos => safe(() => submitCheckin(pos.coords.latitude, pos.coords.longitude, pos.coords.accuracy)),
        () => showCheckinError("Location access is required. Please allow location and try again."),
        { enableHighAccuracy: true, timeout: 10000, maximumAge: 0 }
    );
}

async function submitCheckin(lat, lng, accuracy) {
    const data = await api("POST", "/api/attendance/checkin", { token: qrToken, lat, lng, accuracy }, jwtToken);
    document.getElementById("checkin-time").textContent = data.checkin_at;
    showSection("success-section");
}

function showSection(id) {
    ["login-section", "loading-section", "success-section", "checkin-error-section"]
        .forEach(s => document.getElementById(s).style.display = "none");
    document.getElementById(id).style.display = "block";
}

function showError(msg) {
    const el = document.getElementById("error-msg");
    el.textContent = msg;
    el.style.color = "red";
}

function showCheckinError(msg) {
    document.getElementById("checkin-error-msg").textContent = msg;
    showSection("checkin-error-section");
}

function clearError() { document.getElementById("error-msg").textContent = ""; }

setErrorHandler(showCheckinError);