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
	Verified int64 `json:"verified"`
	Specialty string `json:"specialty"`
	Token string `json:"token"`
}

type StudentRegisterRequest struct {
	StudentID string `json:"student_id"`
	Email string `json:"email"`
	Password string `json:"password"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Specialty string `json:"specialty"`
}

type StudentRegisterResponse struct {
	StudentID string `json:"student_id"`
	Email string `json:"email"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Specialty string `json:"specialty"`
}