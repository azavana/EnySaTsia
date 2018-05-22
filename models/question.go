package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Question definition
type Question struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Sentence string        `json:"sentence"`
	Yes      []string      `json:"yes"`
	No       []string      `json:"no"`
	NoIdea   []string      `json:"noIdea"`
	Voted    bool          `json:"voted"`
	State    string        `json:"state"`
	Session  bson.ObjectId `bson:"session"`
}

// UpdateQuestion -- the update model should only contain an ID and a sentence
type UpdateQuestion struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Sentence string        `json:"sentence"`
}

//Vote --the vote of each User
type Vote struct {
	User     string        `json:"user"`
	Decision string        `json:"decision"`
	Question bson.ObjectId `json:"question"`
}
