let currentSessionID = null;
let qrInterval = null;

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
            td.textContent = val;
            tr.appendChild(td);
        }
        table.appendChild(tr);
    }
    return table;
}

async function populateDropdown(selectId, url, valueFn, labelFn) {
    const items = await api("GET", url);
    const select = document.getElementById(selectId);
    const current = select.value;
    select.innerHTML = '<option value="">-- Select --</option>';
    for (const item of items) {
        const opt = document.createElement("option");
        opt.value = valueFn(item);
        opt.textContent = labelFn(item);
        select.appendChild(opt);
    }
    if (current) select.value = current;
}

function courseLabel(c) {
    return c.course_name + " — " + c.section_date + " " + c.start_time;
}

function refreshAllDropdowns() {
    const ids = ["session-course", "roster-course", "export-course"];
    for (const id of ids) {
        populateDropdown(id, "/api/courses", c => c.course_id, courseLabel);
    }
}

function showMsg(msg) {
    const el = document.getElementById("dashboard-msg");
    el.textContent = msg;
    setTimeout(() => { el.textContent = ""; }, 5000);
}

// === Page load ===
document.addEventListener("DOMContentLoaded", () => {
    loadCourses();
    loadSpecialties();

    document.getElementById("create-course-form").addEventListener("submit", e => {
        e.preventDefault();
        safe(createCourse);
    });
    document.getElementById("create-specialty-form").addEventListener("submit", e => {
        e.preventDefault();
        safe(createSpecialty);
    });
});

// === Tabs ===
function showTab(tabName, btn) {
    document.querySelectorAll(".tab-content").forEach(el => el.style.display = "none");
    document.querySelectorAll(".tab-btn").forEach(el => el.classList.remove("active"));
    document.getElementById("tab-" + tabName).style.display = "block";
    btn.classList.add("active");
    if (["session", "roster", "export"].includes(tabName)) refreshAllDropdowns();
}

// === Courses ===
async function loadCourses() {
    const courses = await api("GET", "/api/courses");
    const list = document.getElementById("course-list");
    list.innerHTML = "";
    if (courses.length === 0) { list.textContent = "No courses yet."; return; }
    list.appendChild(buildTable(
        ["Name", "Day", "Time"],
        courses.map(c => [c.course_name, c.section_date, c.start_time])
    ));
    refreshAllDropdowns();
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
        div.style.cssText = "display:flex;align-items:center;gap:8px;margin-bottom:4px";
        const name = document.createElement("span");
        name.textContent = s.specialty_name;
        const btn = document.createElement("button");
        btn.textContent = "Delete";
        btn.style.cssText = "padding:2px 8px;font-size:12px";
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

    navigator.geolocation.getCurrentPosition(
        pos => safe(async () => {
            const session = await api("POST", "/api/sessions/start", {
                course_id: courseID,
                classroom_lat: pos.coords.latitude,
                classroom_lng: pos.coords.longitude,
            });
            currentSessionID = session.session_id;
            document.getElementById("no-active-session").style.display = "none";
            document.getElementById("active-session").style.display = "block";
            document.getElementById("session-info").textContent = "Session active — " + session.session_date;
            refreshQR();
            qrInterval = setInterval(refreshQR, 13000);
        }),
        () => showMsg("Location access required"),
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
        ["Student ID", "Name", "Email", "Specialty"],
        students.map(s => [s.student_id, s.first_name + " " + s.last_name, s.email, s.specialty || ""])
    ));
}

// === Export ===
function exportSemester() {
    const courseID = document.getElementById("export-course").value;
    if (!courseID) { showMsg("Please select a course"); return; }
    let url = "/api/export/semester/" + courseID;
    const from = document.getElementById("export-from").value;
    const to = document.getElementById("export-to").value;
    if (from && to) url += "?from=" + from + "&to=" + to;
    window.location.href = url;
}

function exportDaily() {
    const date = document.getElementById("export-date").value;
    if (!date) { showMsg("Please select a date"); return; }
    window.location.href = "/api/export/daily/" + date;
}

setErrorHandler(showMsg);