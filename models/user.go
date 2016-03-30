package models

import "time"

type User struct {
	FirstName  string `form:"first_name" binding:"required"`
	LastName string `form:"last_name" binding:"required"`
	ContactNo string `form:"contact_no" binding:"required"`
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Status string `form:"status" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}