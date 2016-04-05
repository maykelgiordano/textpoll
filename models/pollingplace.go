package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"	
)

type PollingPlace struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Place  string `form:"place" binding:"required"`
	Status string `form:"status" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}