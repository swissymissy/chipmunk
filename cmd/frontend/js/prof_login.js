document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("login-form").addEventListener("submit", e => {
        e.preventDefault();
        safe(handleLogin);
    });
});

async function handleLogin() {
    const password = document.getElementById("password").value;
    const data = await api("POST", "/api/auth/professor/login", { password });
    localStorage.setItem("professor_token", data.token);
    window.location.href = "/";
}

function showError(msg) {
    const el = document.getElementById("error-msg");
    el.textContent = msg;
    el.style.color = "red";
}

setErrorHandler(showError);
