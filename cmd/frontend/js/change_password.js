// --- auth bootstrap (mirrors profile.js) ---
const studentToken = localStorage.getItem("student_token");
if (!studentToken) {
    window.location.href = "/student_login.html";
}
setAuthToken(studentToken);
setUnauthorizedHandler(clearStudentAuth);

document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("change-form").addEventListener("submit", e => {
        e.preventDefault();
        submitForm(e.target, handleChange);
    });
});

async function handleChange() {
    clearError();
    const current = document.getElementById("current-password").value;
    const next = document.getElementById("new-password").value;
    const confirm = document.getElementById("confirm-password").value;

    if (next !== confirm) {
        showError("New passwords don't match.");
        return;
    }

    // server verifies the current password and clears the must-change flag.
    await api("PUT", "/api/students/myprofile/password", {
        current_password: current,
        new_password: next,
    });

    // flag is now cleared server-side, so the profile page won't bounce back here.
    window.location.href = "/profile.html";
}

function showError(msg) {
    const el = document.getElementById("error-msg");
    el.textContent = msg;
    el.style.color = "red";
}
function clearError() {
    document.getElementById("error-msg").textContent = "";
}

setErrorHandler(showError);
