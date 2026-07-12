let currentSessionID = null;
let qrInterval = null;
let currentAttendanceSessionID = null;
let attendanceInterval = null;

const professorToken = localStorage.getItem("professor_token");
if (!professorToken) {
    window.location.href = "/prof_login.html";
}
setAuthToken(professorToken);
setUnauthorizedHandler(clearProfessorAuth);

function logout() {
    clearProfessorAuth();
}

function buildTable(headers, rows) {
    const table = document.createElement("table");
    const thead = document.createElement("tr");
    for (const h of headers) {
        const th = document.createElement("th");
        th.textContent = h;
        thead.appendChild(th);
    }
    table.appendChild(thead);
    for (const row of rows) {
        const tr = document.createElement("tr");
        for (const val of row) {
            const td = document.createElement("td");
            if (val instanceof Node) td.appendChild(val);
            else td.textContent = val;
            tr.appendChild(td);
        }
        table.appendChild(tr);
    }
    return table;
}

function fillCourseDropdowns(courses) {
    for (const id of ["session-course", "roster-course", "export-course"]) {
        fillDropdown(id, courses, c => c.course_id, courseLabel);
    }
}

async function refreshAllDropdowns() {
    const courses = await api("GET", "/api/courses");
    fillCourseDropdowns(courses);
}

function showMsg(msg) {
    const el = document.getElementById("dashboard-msg");
    el.textContent = msg;
    setTimeout(() => { el.textContent = ""; }, 5000);
}

// === Page load ===
document.addEventListener("DOMContentLoaded", () => {
    safe(loadCourses);
    safe(loadSpecialties);
    safe(checkForActiveSession);

    document.getElementById("create-course-form").addEventListener("submit", e => {
        e.preventDefault();
        submitForm(e.target, createCourse);
    });
    document.getElementById("create-specialty-form").addEventListener("submit", e => {
        e.preventDefault();
        submitForm(e.target, createSpecialty);
    });
});

function enterActiveSessionUI(session) {
    currentSessionID = session.session_id;
    document.getElementById("no-active-session").style.display = "none";
    document.getElementById("active-session").style.display = "block";
    document.getElementById("session-info").textContent = "Session active — " + session.started_at;
    safe(refreshQR);
    qrInterval = setInterval(() => safe(refreshQR), 13000);
}

async function checkForActiveSession() {
    const sessions = await api("GET", "/api/sessions/active");
    if (sessions.length === 0) return;
    if (sessions.length > 1) console.warn("multiple active sessions; resuming first", sessions);
    enterActiveSessionUI(sessions[0]);
}


// === Tabs ===
function showTab(tabName, btn) {
    document.querySelectorAll(".tab-content").forEach(el => el.style.display = "none");
    document.querySelectorAll(".tab-btn").forEach(el => el.classList.remove("active"));
    document.getElementById("tab-" + tabName).style.display = "block";
    btn.classList.add("active");

    // stop attendance polling if leaving
    if (attendanceInterval) {
        clearInterval(attendanceInterval);
        attendanceInterval = null;
    }

    if (["session", "roster", "export"].includes(tabName)) safe(refreshAllDropdowns);

    if (tabName === "attendance") {
        loadAttendance();
        attendanceInterval = setInterval(() => {
            if (currentAttendanceSessionID) loadAttendanceRoster(currentAttendanceSessionID);
        }, 5000);
    }

    if (tabName === "settings") renderIndividualResets();
}

// === Courses ===
async function loadCourses() {
    const courses = await api("GET", "/api/courses");
    const list = document.getElementById("course-list");
    list.innerHTML = "";
    if (courses.length === 0) { list.textContent = "No courses yet."; return; }

    const rows = courses.map(c => {
        const btn = document.createElement("button");
        btn.textContent = "Delete";
        btn.onclick = () => safe(() => deleteCourse(c.course_id, c.course_name));
        return [c.course_name, c.section_date, c.start_time, btn];
    });

    list.appendChild(buildTable(
        ["Name", "Day", "Time", "Action"],
        rows,
    ));
    fillCourseDropdowns(courses);
}

async function deleteCourse(courseID, courseName) {
    if (!confirm(`Delete course "${courseName}"?\n\nThis will also delete all sessions, enrollments, and attendance records for this course.`)) {
        return;
    }
    await api("DELETE", "/api/courses/" + courseID);
    showMsg("Course deleted");
    loadCourses();
}

async function createCourse() {
    await api("POST", "/api/courses", {
        name: document.getElementById("course-name").value.trim(),
        section_date: document.getElementById("course-day").value.trim(),
        start_time: document.getElementById("course-time").value.trim(),
    });
    document.getElementById("create-course-form").reset();
    loadCourses();
    showMsg("Course created!");
}

// === Specialties ===
async function loadSpecialties() {
    const specialties = await api("GET", "/api/specialties");
    const list = document.getElementById("specialty-list");
    list.innerHTML = "";
    if (specialties.length === 0) { list.textContent = "No specialties yet."; return; }
    for (const s of specialties) {
        const div = document.createElement("div");
        div.className = "specialty-row";
        const name = document.createElement("span");
        name.textContent = s.specialty_name;
        const btn = document.createElement("button");
        btn.textContent = "Delete";
        btn.onclick = () => safe(async () => {
            await api("DELETE", "/api/specialties/" + s.id);
            loadSpecialties();
        });
        div.append(name, btn);
        list.appendChild(div);
    }
}

async function createSpecialty() {
    await api("POST", "/api/specialties", {
        specialty_name: document.getElementById("specialty-name").value.trim(),
    });
    document.getElementById("create-specialty-form").reset();
    loadSpecialties();
    showMsg("Specialty added!");
}

// === Session ===
function startSession() {
    const courseID = document.getElementById("session-course").value;
    if (!courseID) { showMsg("Please select a course"); return; }
    if (!navigator.geolocation) { showMsg("Location not supported"); return; }

    const btn = document.getElementById("start-session-btn");
    if (btn) btn.disabled = true;

    navigator.geolocation.getCurrentPosition(
        pos => safe(async () => {
            try {
                const session = await api("POST", "/api/sessions/start", {
                    course_id: courseID,
                    classroom_lat: pos.coords.latitude,
                    classroom_lng: pos.coords.longitude,
                });
                enterActiveSessionUI(session);
            } finally {
                if (btn) btn.disabled = false;
            }
        }),
        () => { showMsg("Location access required"); if (btn) btn.disabled = false; },
        { enableHighAccuracy: true, timeout: 10000 }
    );
}

async function refreshQR() {
    const data = await api("GET", "/api/sessions/" + currentSessionID + "/qr");
    const container = document.getElementById("qr-container");
    container.innerHTML = "";
    new QRCode(container, { text: data.checkin_url, width: 300, height: 300 });
}

async function closeSession() {
    await api("PUT", "/api/sessions/close", { session_id: currentSessionID });
    clearInterval(qrInterval);
    qrInterval = null;
    currentSessionID = null;
    document.getElementById("active-session").style.display = "none";
    document.getElementById("no-active-session").style.display = "block";
    document.getElementById("qr-container").innerHTML = "";
    showMsg("Session closed");
}

// === Roster ===
async function loadRoster() {
    const courseID = document.getElementById("roster-course").value;
    if (!courseID) { showMsg("Please select a course"); return; }
    const students = await api("GET", "/api/roster/" + courseID);
    const list = document.getElementById("roster-list");
    list.innerHTML = "";
    if (students.length === 0) { list.textContent = "No students enrolled."; return; }
    list.appendChild(buildTable(
        ["Student ID", "Name", "Email", "Specialty", "Action"],
        students.map(s => {
            // Remove button — unenrolls the student from this course only.
            // s.id is the internal student id that enrollments key on
            // (not s.student_id, which is the school ID).
            const remove = document.createElement("button");
            remove.textContent = "Remove";
            remove.onclick = () => removeStudentFromCourse(courseID, s.id, s.first_name + " " + s.last_name);
            return [s.student_id, s.first_name + " " + s.last_name, s.email, s.specialty || "", remove];
        })
    ));
}

// unenroll a student from a course (does not delete their account).
async function removeStudentFromCourse(courseID, studentID, name) {
    if (!confirm(`Remove ${name} from this course? This only unenrolls them — their account is kept.`)) return;
    await safe(async () => {
        await api("DELETE", "/api/roster/" + courseID + "/students/" + studentID);
        showMsg("Student removed from course");
        await loadRoster();
    });
}

// === Export ===
// Exports use fetch + blob (not navigation) so the JWT travels in the
// Authorization header. window.location.href would skip the header entirely.
function exportSemester() {
    const courseID = document.getElementById("export-course").value;
    if (!courseID) { showMsg("Please select a course"); return; }
    let url = "/api/export/semester/" + courseID;
    const from = document.getElementById("export-from").value;
    const to = document.getElementById("export-to").value;
    if (from && to) url += "?from=" + from + "&to=" + to;
    safe(() => downloadFile(url, "semester_report.xlsx"));
}

function exportDaily() {
    const date = document.getElementById("export-date").value;
    if (!date) { showMsg("Please select a date"); return; }
    safe(() => downloadFile("/api/export/daily/" + date, date + "_report.xlsx"));
}

async function downloadFile(url, fallbackName) {
    const token = localStorage.getItem("professor_token");
    const res = await fetch(url, {
        headers: token ? { Authorization: "Bearer " + token } : {},
    });
    if (res.status === 401) {
        clearProfessorAuth();
        return;
    }
    if (!res.ok) throw new Error("Download failed (" + res.status + ")");

    // server already sets Content-Disposition; read the filename from there
    let filename = fallbackName;
    const cd = res.headers.get("Content-Disposition");
    if (cd) {
        const m = cd.match(/filename="?([^"]+)"?/);
        if (m) filename = m[1];
    }

    const blob = await res.blob();
    const objectUrl = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = objectUrl;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(objectUrl);
}

setErrorHandler(showMsg);

// runs when Attendance tab is opened
// fetches active sessions + course names in parallel
// populate the pickers, then load the rosters
async function loadAttendance() {
    try {
        const [sessions, courses] = await Promise.all([
            api("GET", "/api/sessions/active"),
            api("GET", "/api/courses"),
        ]);

        const empty = document.getElementById("attendance-empty");
        const controls = document.getElementById("attendance-controls");
        const list = document.getElementById("attendance-list");

        if (sessions.length === 0) {
            empty.style.display = "block";
            controls.style.display = "none";
            list.innerHTML = "";
            renderFlagGroups([]);
            currentAttendanceSessionID = null;
            return;
        }
        empty.style.display = "none";

        const courseName = new Map(courses.map(c => [c.course_id, c.course_name]));

        const select = document.getElementById("attendance-session");
        select.innerHTML = "";
        for (const s of sessions) {
            const opt = document.createElement("option");
            opt.value = s.session_id;
            opt.textContent = (courseName.get(s.course_id) || s.course_id) + " — " + s.session_date;
            select.appendChild(opt);
        }

        controls.style.display = sessions.length > 1 ? "block" : "none";

        // keep prior selection if still active, else default to first
        const prev = currentAttendanceSessionID;
        const ids = new Set(sessions.map(s => s.session_id));
        currentAttendanceSessionID = (prev && ids.has(prev)) ? prev : sessions[0].session_id;
        select.value = String(currentAttendanceSessionID);

        await loadAttendanceRoster(currentAttendanceSessionID);
    } catch (err) {
        const msgEl = document.getElementById("attendance-msg");
        msgEl.style.color = "";
        msgEl.textContent = err.message;
    }
}

// fired by the picker when prof changes selection
function onAttendanceSessionChange() {
    currentAttendanceSessionID = parseInt(document.getElementById("attendance-session").value);
    loadAttendanceRoster(currentAttendanceSessionID);
}

// fetches the attendance roster for a session and renders the table.
// "Mark Present" button only appears for absent rows.
// uses local try/catch (not safe()) so errors stay in attendance-msg
// instead of replacing the whole tab via the global error handler.
async function loadAttendanceRoster(sessionID) {
    if (!sessionID) return;
    const msgEl = document.getElementById("attendance-msg");
    try {
        const data = await api("GET", "/api/attendance/" + sessionID);
        const rows = data.roster || [];
        const flagGroups = data.flag_groups || [];

        const list = document.getElementById("attendance-list");
        list.innerHTML = "";
        if (rows.length === 0) {
            list.textContent = "No students enrolled.";
            renderFlagGroups([]); // clear any leftover card
            return;
        }

        const tableRows = rows.map(r => {
            // Name cell - wrap in a span so we can add a badge
            const nameCell = document.createElement("span");
            nameCell.textContent = r.first_name + " " + r.last_name;
            if (r.flagged) {
                const badge = document.createElement("span");
                badge.textContent = " ⚠";
                badge.title = "Shares device with another student in this session";
                badge.className = "flag-badge";
                nameCell.appendChild(badge);
            }

            // Action cell - mark present (if absent) or mark absent (if present)
            const action = document.createElement("button");
            if (r.status === "absent") {
                action.textContent = "Mark Present";
                action.onclick = () => markPresent(r.student_id, r.session_id);
            } else {
                action.textContent = "Mark Absent";
                action.onclick = () => markAbsent(r.student_id, r.session_id);
            }

            return [
                r.student_school_id,
                nameCell,
                r.status,
                r.checkin_at || "—",
                action,
            ];
        });

        list.appendChild(buildTable(
            ["Student ID", "Name", "Status", "Checked-in at", "Action"],
            tableRows,
        ));

        renderFlagGroups(flagGroups);
    } catch (err) {
        msgEl.style.color = "";
        msgEl.textContent = err.message;
    }
}

function renderFlagGroups(groups) {
    const container = document.getElementById("attendance-flags");
    container.innerHTML = "";
    if (groups.length === 0) return;

    const heading = document.createElement("h3");
    heading.textContent = "Suspicious Activity";
    container.appendChild(heading);

    for (const g of groups) {
        const card = document.createElement("div");
        card.className = "flag-card";

        const title = document.createElement("div");
        title.className = "flag-card-title";
        title.textContent = `Shared device — ${g.students.length} students`;
        card.appendChild(title);

        const ul = document.createElement("ul");
        for (const s of g.students) {
            const li = document.createElement("li");
            li.textContent = `${s.first_name} ${s.last_name} (${s.student_school_id}) — checked in ${s.check_in_at}`;
            ul.appendChild(li);
        }
        card.appendChild(ul);
        container.appendChild(card);
    }
}

// flip a student's attendance status via the prof override endpoints.
// shared by markPresent and markAbsent — they only differ in the URL.
async function flipAttendance(path, studentID, sessionID) {
    const msgEl = document.getElementById("attendance-msg");
    msgEl.style.color = "";
    msgEl.textContent = "";
    try {
        await api("PUT", path, { student_id: studentID, session_id: sessionID });
        await loadAttendanceRoster(currentAttendanceSessionID);
    } catch (err) {
        msgEl.style.color = "red";
        msgEl.textContent = err.message;
    }
}

const markPresent = (s, ses) => flipAttendance("/api/attendance/override", s, ses);
const markAbsent  = (s, ses) => flipAttendance("/api/attendance/override/absent", s, ses);

// === Settings: destructive resets ===
// Each row toggles its button enabled only when the matching input contains "RESET".
// onResetConfirmInput is wired via inline oninput=...; rendered rows wire it dynamically.

const INDIVIDUAL_RESETS = [
    { key: "records",     label: "Attendance Records",   path: "/api/reset/records" },
    { key: "sessions",    label: "Attendance Sessions",  path: "/api/reset/sessions" },
    { key: "enrollments", label: "Enrollments",          path: "/api/reset/enrollments" },
    { key: "courses",     label: "Courses",              path: "/api/reset/courses" },
    { key: "students",    label: "Students",             path: "/api/reset/students" },
    { key: "specialties", label: "Specialties",          path: "/api/reset/specialties" },
];

function onResetConfirmInput(inputId, btnId) {
    const input = document.getElementById(inputId);
    const btn = document.getElementById(btnId);
    btn.disabled = input.value.trim() !== "RESET";
}

async function resetAll() {
    await api("DELETE", "/api/reset/all");
    document.getElementById("reset-all-confirm").value = "";
    document.getElementById("reset-all-btn").disabled = true;
    showMsg("Database reset for new semester");
    // refresh views that depend on now-cleared data
    safe(loadCourses);
    safe(loadSpecialties);
}

async function resetTable(path, inputId, btnId) {
    await api("DELETE", path);
    document.getElementById(inputId).value = "";
    document.getElementById(btnId).disabled = true;
    showMsg("Reset complete");
    safe(loadCourses);
    safe(loadSpecialties);
}

function renderIndividualResets() {
    const container = document.getElementById("individual-resets");
    if (!container || container.dataset.rendered) return;
    container.dataset.rendered = "1";

    for (const r of INDIVIDUAL_RESETS) {
        const inputId = `reset-${r.key}-confirm`;
        const btnId = `reset-${r.key}-btn`;

        const row = document.createElement("div");
        row.className = "reset-row";

        const label = document.createElement("label");
        label.htmlFor = inputId;
        label.textContent = r.label;

        const input = document.createElement("input");
        input.type = "text";
        input.id = inputId;
        input.placeholder = "Type RESET";
        input.autocomplete = "off";
        input.oninput = () => onResetConfirmInput(inputId, btnId);

        const btn = document.createElement("button");
        btn.id = btnId;
        btn.className = "danger-btn";
        btn.textContent = "Reset";
        btn.disabled = true;
        btn.onclick = () => safe(() => resetTable(r.path, inputId, btnId));

        row.append(label, input, btn);
        container.appendChild(row);
    }
}
