document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("login-form").addEventListener("submit", e => {
        e.preventDefault();
        submitForm(e.target, handleLogin);
    });
});

async function handleLogin() {
    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value;

    const data = await api("POST", "/api/auth/login", { email, password });
    // store under a student-specific key so it never collides with the
    // professor session.
    localStorage.setItem("student_token", data.token);
    window.location.href = "/profile.html";
}

function showError(msg) {
    const el = document.getElementById("error-msg");
    el.textContent = msg;
    el.style.color = "red";
}

setErrorHandler(showError);
