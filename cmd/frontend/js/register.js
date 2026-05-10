document.addEventListener("DOMContentLoaded", () => {
    safe(loadCourses);
    safe(loadSpecialties);
    document.getElementById("register-form").addEventListener("submit", e => {
        e.preventDefault();
        submitForm(e.target, handleRegister);
    });
});

async function loadCourses() {
    const courses = await api("GET", "/api/courses");
    fillDropdown("course", courses, c => c.course_id, courseLabel, "-- Select a course --");
}

async function loadSpecialties() {
    const specialties = await api("GET", "/api/specialties");
    fillDropdown("specialty", specialties, s => s.specialty_name, s => s.specialty_name, "-- Select your major --");
}

async function handleRegister() {
    clearMessages();
    const studentID = document.getElementById("student_id").value.trim();
    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value;
    const firstName = document.getElementById("first_name").value.trim();
    const lastName = document.getElementById("last_name").value.trim();
    const specialty = document.getElementById("specialty").value;
    const courseID = document.getElementById("course").value;

    if (password.length < 8) { showError("Password should be at least 8 characters"); return; }
    if (!courseID) { showError("Please select a course"); return; }

    const device_fingerprint = await getDeviceFingerprint();

    await api("POST", "/api/auth/register", {
        student_id: studentID, email, password,
        first_name: firstName, last_name: lastName, specialty,
        device_fingerprint,
    });

    const loginData = await api("POST", "/api/auth/login", { email, password });
    await api("POST", "/api/enrollment", { course_id: courseID }, loginData.token);

    showSuccess("Registration complete! You are enrolled in the course.");
}

function showError(msg) {
    const el = document.getElementById("error-msg");
    el.textContent = msg;
    el.style.color = "red";
}
function showSuccess(msg) {
    const el = document.getElementById("success-msg");
    el.textContent = msg;
    el.style.color = "green";
}
function clearMessages() {
    document.getElementById("error-msg").textContent = "";
    document.getElementById("success-msg").textContent = "";
}

setErrorHandler(showError);