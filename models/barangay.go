package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"	
)

type Barangay struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	BarangayName  string `form:"brgy_name" binding:"required"`
	Population int `form:"population_count" binding:"required"`
	Status string `form:"status" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}