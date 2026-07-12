// --- auth bootstrap (mirrors dashboard.js) ---
const studentToken = localStorage.getItem("student_token");
if (!studentToken) {
    window.location.href = "/student_login.html";
}
setAuthToken(studentToken);
setUnauthorizedHandler(clearStudentAuth);

function logout() {
    clearStudentAuth();
}

// top-level errors (e.g. failed initial load) land in the page-level banner
function showProfileMsg(msg) {
    const el = document.getElementById("profile-msg");
    el.textContent = msg;
    el.style.color = "red";
}
setErrorHandler(showProfileMsg);

document.addEventListener("DOMContentLoaded", () => {
    safe(loadProfile);
    safe(loadMyCourses);
    setupEditableFields();
});

// --- profile fields ---
function fillProfile(p) {
    document.getElementById("school-id").value = p.student_school_id;
    document.getElementById("email").value = p.email;
    document.getElementById("first-name").value = p.first_name;
    document.getElementById("last-name").value = p.last_name;
    document.getElementById("specialty-display").textContent = p.specialty || "—";
}

async function loadProfile() {
    const p = await api("GET", "/api/students/myprofile");
    fillProfile(p);
}

// per-field message helper: green on success, red on error
function fieldMsg(id, text, ok) {
    const el = document.getElementById(id);
    el.style.color = ok ? "green" : "red";
    el.textContent = text;
}

// wire one editable field group. inputs start read-only; the student must
// click Edit to unlock them, then Save (persist) or Cancel (revert). This
// keeps fields locked by default so they can't be changed by accident.
// `save` receives the trimmed input values and returns the update promise.
function setupEditableField({ inputs, prefix, msg, save }) {
    const inputEls = inputs.map(id => document.getElementById(id));
    const editBtn = document.getElementById(prefix + "-edit");
    const saveBtn = document.getElementById(prefix + "-save");
    const cancelBtn = document.getElementById(prefix + "-cancel");
    let original = [];

    function setEditing(on) {
        inputEls.forEach(el => { el.readOnly = !on; });
        editBtn.style.display = on ? "none" : "";
        saveBtn.style.display = on ? "" : "none";
        cancelBtn.style.display = on ? "" : "none";
    }

    editBtn.onclick = () => {
        original = inputEls.map(el => el.value);   // remember, so Cancel can revert
        fieldMsg(msg, "", true);
        setEditing(true);
        inputEls[0].focus();
    };

    cancelBtn.onclick = () => {
        inputEls.forEach((el, i) => { el.value = original[i]; });
        fieldMsg(msg, "", true);
        setEditing(false);
    };

    saveBtn.onclick = async () => {
        fieldMsg(msg, "", true);
        saveBtn.disabled = true;
        try {
            const updated = await save(inputEls.map(el => el.value.trim()));
            if (updated) fillProfile(updated);
            fieldMsg(msg, "Saved!", true);
            setEditing(false);
        } catch (err) {
            fieldMsg(msg, err.message, false);   // stay in edit mode so they can fix it
        } finally {
            saveBtn.disabled = false;
        }
    };
}

function setupEditableFields() {
    setupEditableField({
        inputs: ["school-id"], prefix: "school-id", msg: "school-id-msg",
        save: ([student_school_id]) =>
            api("PUT", "/api/students/myprofile/schoolid", { student_school_id }),
    });
    setupEditableField({
        inputs: ["email"], prefix: "email", msg: "email-msg",
        save: ([email]) =>
            api("PUT", "/api/students/myprofile/email", { email }),
    });
    setupEditableField({
        inputs: ["first-name", "last-name"], prefix: "name", msg: "name-msg",
        save: ([first_name, last_name]) =>
            api("PUT", "/api/students/myprofile/name", { first_name, last_name }),
    });
}

// --- courses ---
async function loadMyCourses() {
    const [enrolled, all] = await Promise.all([
        api("GET", "/api/enrollments"),
        api("GET", "/api/courses"),
    ]);

    const list = document.getElementById("enrolled-list");
    list.innerHTML = "";
    for (const c of enrolled) {
        const li = document.createElement("li");
        li.className = "course-row";

        const span = document.createElement("span");
        span.textContent = courseLabel(c);
        li.appendChild(span);

        const btn = document.createElement("button");
        btn.textContent = "Remove";
        btn.onclick = () => removeCourse(c.course_id);
        li.appendChild(btn);

        list.appendChild(li);
    }

    const enrolledIds = new Set(enrolled.map(c => c.course_id));
    const available = all.filter(c => !enrolledIds.has(c.course_id));
    fillDropdown("add-course", available, c => c.course_id, courseLabel, "-- Select a course --");
}

async function addCourse() {
    const courseID = document.getElementById("add-course").value;
    if (!courseID) return;
    const btn = document.getElementById("add-course-btn");
    const msgEl = document.getElementById("add-course-msg");
    msgEl.textContent = "";
    if (btn) btn.disabled = true;
    try {
        await api("POST", "/api/enrollment", { course_id: courseID });
        msgEl.style.color = "green";
        msgEl.textContent = "Course added!";
        await loadMyCourses();
    } catch (err) {
        msgEl.style.color = "red";
        msgEl.textContent = err.message;
    } finally {
        if (btn) btn.disabled = false;
    }
}

async function removeCourse(courseID) {
    const msgEl = document.getElementById("add-course-msg");
    msgEl.textContent = "";
    try {
        await api("DELETE", "/api/students/myprofile/courses/" + courseID);
        await loadMyCourses();
    } catch (err) {
        msgEl.style.color = "red";
        msgEl.textContent = err.message;
    }
}
