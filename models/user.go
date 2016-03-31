package models

import "time"

type User struct {
	FirstName  string `form:"first_name" binding:"required"`
	LastName string `form:"last_name" binding:"required"`
	ContactNo string `form:"contact_no" binding:"required"`
	Email string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
	Status string `form:"status" binding:"required"`
	UserRole string `form:"user_role" biding:"required"`
	IsDefaultPassword bool
	CreatedAt time.Time
	UpdatedAt time.Time
}