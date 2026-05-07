let jwtToken = null;
let qrToken = null;
let firstName = "";

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
    firstName = data.first_name;

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
    document.getElementById("success-greeting").textContent = "Hi " + firstName + ", you've checked in!";
    document.getElementById("checkin-time").textContent = "Checked in at " + data.check_in_at;
    showSection("success-section");
    loadMyCourses();
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

async function loadMyCourses() {
    try {
        const [enrolled, all] = await Promise.all([
            api("GET", "/api/enrollments", null, jwtToken),
            api("GET", "/api/courses"),
        ]);

        const list = document.getElementById("enrolled-list");
        list.innerHTML = "";
        for (const c of enrolled) {
            const li = document.createElement("li");
            li.textContent = courseLabel(c);
            list.appendChild(li);
        }

        const enrolledIds = new Set(enrolled.map(c => c.course_id));
        const available = all.filter(c => !enrolledIds.has(c.course_id));
        fillDropdown("add-course", available, c => c.course_id, courseLabel, "-- Select a course --");
    } catch (err) {
        document.getElementById("enrolled-list").innerHTML =
            "<li>Couldn't load courses: " + err.message + "</li>";
    }
}

async function addCourse() {
    const courseID = document.getElementById("add-course").value;
    if (!courseID) return;
    const msgEl = document.getElementById("add-course-msg");
    msgEl.textContent = "";
    try {
        await api("POST", "/api/enrollment", { course_id: courseID }, jwtToken);
        msgEl.style.color = "green";
        msgEl.textContent = "Course added!";
        await loadMyCourses();
    } catch (err) {
        msgEl.style.color = "red";
        msgEl.textContent = err.message;
    }
}
