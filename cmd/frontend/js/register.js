document.addEventListener("DOMContentLoaded", () => {
    safe(loadCourses);
    safe(loadSpecialties);
    document.getElementById("register-form").addEventListener("submit", e => {
        e.preventDefault();
        safe(handleRegister);
    });
});

async function loadCourses() {
    const courses = await api("GET", "/api/courses");
    const select = document.getElementById("course");
    select.innerHTML = '<option value="">-- Select a course --</option>';
    for (const c of courses) {
        const opt = document.createElement("option");
        opt.value = c.course_id;
        opt.textContent = c.course_name + " — " + c.section_date + " " + c.start_time;
        select.appendChild(opt);
    }
}

async function loadSpecialties() {
    const specialties = await api("GET", "/api/specialties");
    const select = document.getElementById("specialty");
    select.innerHTML = '<option value="">-- Select your major --</option>';
    for (const s of specialties) {
        const opt = document.createElement("option");
        opt.value = s.specialty_name;
        opt.textContent = s.specialty_name;
        select.appendChild(opt);
    }
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

    await api("POST", "/api/auth/register", {
        student_id: studentID, email, password,
        first_name: firstName, last_name: lastName, specialty,
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