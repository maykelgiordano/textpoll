package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"	
)

type PollingPlace struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Place  string `json:"place" form:"place" binding:"required"`
	Status string `json:"status" form:"status" binding:"required"`
	BarangayId string `form:"barangay_id" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}