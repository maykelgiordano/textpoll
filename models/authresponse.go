package models

type CurrentUser struct {
	FirstName  string `json:"first_name"`
	LastName string `json:"last_name"`
	ContactNo string `json:"contact_no"`
	Email string `json:"email"`
	UserRole string `json:"user_role"`
	IsDefaultPassword bool `json:"is_default_password"`
	Token string `json:"token"`
}