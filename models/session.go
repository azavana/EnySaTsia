package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Session definition
type Session struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Name   string        `json:"name"`
	Voters int           `json:"voters"`
	State  string        `json:"State"`
}
