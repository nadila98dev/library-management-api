package model

type Users struct {
	ID string `json:"id"`
	StudentID string `json:"student_id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role string `json:"role"`
	Password string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}