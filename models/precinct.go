package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"	
)

type Precinct struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	PrecinctNo  string `json:"precinct_no" form:"precinct_no" binding:"required"`
	RegisteredVoters int `json:"registered_voters" form:"registered_voters" binding:"required"`
	PollingPlaceId string `json:"polling_place_id" form:"polling_place_id" binding:"required"`
	Status string `json:"status" form:"status" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}