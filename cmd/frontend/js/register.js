// runs when the page loads
document.addEventListener("DOMContentLoaded", function () {
    loadCourses();

    document.getElementById("register-form").addEventListener("submit", function (e) {
        e.preventDefault();
        handleResgiter();
    });
});

// fetch courses from API and populate the dropdown
async function loadCourses() {
    try {
        const res = await fetch("api/courses");
        const courses = await res.json();

        const select = document.getElementById("course");
        select.innerHTML = '<option value="">-- Select a course --</option>';

        for (const course of courses) {
            const option = document.createElement("option");
            option.value = course.course_id;
            option.textContent = course.course_name + "-" + course.section + " " + course.time;
            select.appendChild(option);
        }
    } catch (err) {
        showError("failed to load courses")
    }
}

// handle the registration submit
async function handleResgiter() {
    clearMessages();

    const studentID = document.getElementById("student_id").value.trim();
    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value.trim();
    const firstName = document.getElementById("first_name").value.trim();
    const lastName = document.getElementById("last_name").value.trim();
    const specialty = document.getElementById("specialty").value.trim();
    const courseID = document.getElementById("course").value;

    if (password.length < 8) {
        showError("Password should be at least 8 characters")
        return;
    }

    if (!courseID) {
        showError("Please select a course to register");
        return;
    }

    // register the student
    try {
        const registerRes = await fetch("/api/auth/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                student_id: studentID,
                email: email,
                password: password,
                first_name: firstName,
                last_name: lastName,
                specialty: specialty,
            }),
        });

        if (!registerRes.ok) {
            const err = await registerRes.json();
            showError(err.error || "Registration failed");
            return;
        }

        const student = await registerRes.json();

        // log in to get JWT
        const loginRes = await fetch("/api/auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email: email,
                password: password,
            }),
        });

        if (!loginRes.ok) {
            showError("Account created but login failed. Try logging in manually.");
            return;
        }

        const loginData = await loginRes.json();
        const token = loginData.token;

        // enroll in course using JWT
        const enrollRes = await fetch("/api/enrollment", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + token,
            },
            body: JSON.stringify({
                course_id: courseID,
            }),
        });

        if (!enrollRes.ok) {
            showError("Account created but enrollment failed. Contact your professor.");
            return;
        }

        showSuccess("Registration complete! You are enrolled in the course.");
    } catch (err) {
        showError("Something went wrong. Please try again.");
    }
}

function showError(msg) {
    document.getElementById("error-msg").textContent = msg;
    document.getElementById("error-msg").style.color = "red";
}

function showSuccess(msg) {
    document.getElementById("success-msg").textContent = msg;
    document.getElementById("success-msg").style.color = "green";
}

function clearMessages() {
    document.getElementById("error-msg").textContent = "";
    document.getElementById("success-msg").textContent = "";
}