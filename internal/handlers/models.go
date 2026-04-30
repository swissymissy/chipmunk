package handlers

type StudentLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type StudentLoginResponse struct {
	StudentID string `json:"student_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Verified  int64  `json:"verified"`
	Specialty string `json:"specialty"`
	Token     string `json:"token"`
}

type StudentRegisterRequest struct {
	StudentID string `json:"student_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Specialty string `json:"specialty"`
}
type StudentRegisterResponse struct {
	StudentID string `json:"student_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Specialty string `json:"specialty"`
}

type NewEnrollmentRequest struct {
	CourseID string `json:"course_id"`
}

type NewCourseRequest struct {
	Name    string `json:"name"`
	Section string `json:"section_date"`
	Time    string `json:"start_time"`
}
type NewCourseResponse struct {
	ID      string `json:"course_id"`
	Name    string `json:"name"`
	Section string `json:"section_date"`
	Time    string `json:"start_time"`
}

type Course struct {
	ID         string `json:"course_id"`
	CourseName string `json:"course_name"`
	Section    string `json:"section"`
	Time       string `json:"time"`
}

type StartSessionRequest struct {
	CourseID     string  `json:"course_id"`
	ClassroomLat float64 `json:"classroom_lat"`
	ClassroomLng float64 `json:"classroom_lng"`
}
type StartSessionResponse struct {
	SessionID   int64  `json:"session_id"`
	CourseID    string `json:"course_id"`
	SessionDate string `json:"session_date"`
	Status      string `json:"status"`
	StartedAt   string `json:"started_at"`
}

type CloseSessionRequest struct {
	SessionID int64 `json:"session_id"`
}
type CloseSessionResponse struct {
	SessionID   int64  `json:"session_id"`
	CourseID    string `json:"course_id"`
	SessionDate string `json:"session_date"`
	Status      string `json:"status"`
	StartedAt   string `json:"started_at"`
	EndedAt     string `json:"ended_at"`
}
