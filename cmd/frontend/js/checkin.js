let jwtToken = null;
let qrToken = null;

document.addEventListener("DOMContentLoaded", function () {
    // grab QR token from URL
    const params = new URLSearchParams(window.location.search);
    qrToken = params.get("t");

    if (!qrToken) {
        showCheckinError("No check-in code found. Please scan the QR code.");
        return;
    }

    // set up login form
    document.getElementById("login-form").addEventListener("submit", function (e) {
        e.preventDefault();
        handleLogin();
    });
});

async function handleLogin() {
    clearError();

    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value;

    try {
        const res = await fetch("/api/auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email: email,
                password: password
            }),
        });

        if (!res.ok) {
            const err = await res.json();
            showError(err.error || "Login failed");
            return;
        }

        const data = await res.json();
        jwtToken = data.token;
        
        // start loading state with student's name
        document.getElementById("greeting").textContent = "Hi" + data.first_name + "!";
        showSection("loading-section");

        // start check-in process
        startCheckIn();
    } catch (err) {
        showError("Something went wrong. Please try again.")
    }
}

function startCheckIn() {
    // request GPS 
    if (!navigator.geolocation) {
        showCheckinError("Your brower does not support location services.")
    }

    navigator.geolocation.getCurrentPosition(
        function (pos) {
            // got location - send checkin req
            submitCheckin(pos.coords.latitude, pos.coords.longitude, pos.coords.accuracy);
        },
        function (err) {
            showCheckinError("Location access is required to checkin. Please allow location and try again.");
        },
        {
            enableHighAccuracy: true,
            timeout: 10000,
            maximumAge: 0,
        }
    );
}

async function submitCheckin(lat, lng, accuracy) {
    try {
        const res = await fetch("/api/attendance/checkin", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + jwtToken,
            },
            body: JSON.stringify({
                token: qrToken,
                let: lat,
                lng: lng,
                accuracy: accuracy,
            }),
        });

        if (!res.ok) {
            const err = await res.json();
            showCheckinError(err.error || "Check-in failed");
            return;
        }

        const data = await res.json();
        document.getElementById("Checkin-time").textContent = data.checkin_at;
        showSection("success-section");
    } catch (err) {
        showCheckinError("Something went wrong. Please try again.");
    }
}

// UI helpers - show/hide sections
function showSection(sectionID) {
    document.getElementById("login-section").style.display = "none";
    document.getElementById("loading-section").style.display = "none";
    document.getElementById("success-section").style.display = "none";
    document.getElementById("checkin-error-section").style.display = "none";
    document.getElementById(sectionID).style.display = "block";
}

function showError(msg) {
    document.getElementById("error-msg").textContent = msg;
    document.getElementById("error-msg").style.color = "red";
}

function showCheckinError(msg) {
    document.getElementById("checkin-error-msg").textContent = msg;
    showSection("checkin-error-section");
}

function clearError() {
    document.getElementById("error-msg").textContent = "";
}