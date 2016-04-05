package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"	
)

type Barangay struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	BarangayName  string `json:"barangay_name" form:"brgy_name" binding:"required"`
	Status string `json:"status" form:"status" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_ad"`
}