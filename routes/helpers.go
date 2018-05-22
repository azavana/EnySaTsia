package routes

import (
	db "voting/database"
	"voting/models"

	"gopkg.in/mgo.v2/bson"
)

// IsSessionModifiable --return true if session is open, false otherwise
func IsSessionModifiable(sessionID bson.ObjectId) bool {
	isModifiable := false
	//connection to db
	session := db.Connect()
	defer session.Close()
	c := session.DB(db.Database).C("sessions")
	s := models.Session{}

	//find the session with sessionId
	err := c.Find(bson.M{"_id": sessionID}).One(&s)

	//if there are no error
	if err == nil {
		if s.State == "created" {
			isModifiable = true
		} else {
			isModifiable = false
		}
	}
	return isModifiable
}

// SessionExistsAndNotClosed -- return true if session exist and state is not closed
func SessionExistsAndNotClosed(sessionID bson.ObjectId) bool {
	result := false
	//connection to db
	session := db.Connect()
	defer session.Close()
	c := session.DB(db.Database).C("sessions")
	s := models.Session{}

	//find the session with sessionId
	err := c.Find(bson.M{"_id": sessionID}).One(&s)

	//if there are no error
	if err == nil {
		if s.State == "close" {
			result = false
		} else {
			result = true
		}
	}
	return result
}

// IsQuestionModifiable --return true if voted is false
func IsQuestionModifiable(questionID bson.ObjectId) bool {
	result := false
	session := db.Connect()
	defer session.Close()
	c := session.DB(db.Database).C("questions")
	question := models.Question{}

	//find the question with questionId
	err := c.Find(bson.M{"_id": questionID}).One(&question)

	if err == nil {
		if question.Voted && question.State != "open" {
			result = false
		} else {
			result = true
		}
	}
	return result
}

// CanOpen -- can open the vote?
func CanOpen(question models.Question) bool {
	result := false
	if question.State == "created" {
		result = true
	} else {
		result = false
	}
	return result
}

// CanClose -- can close the vote?
func CanClose(question models.Question) bool {
	result := false
	if question.State == "open" {
		result = true
	} else {
		result = false
	}
	return result
}

// IsVoteOpen -- check if vote for a question is open
func IsVoteOpen(question models.Question) bool {
	result := false
	if question.State == "open" {
		result = true
	} else {
		result = false
	}
	return result
}

//HasVoted -- Chech if an user has already voded the question
func HasVoted(question models.Question, userID string) bool {
	result := false
	votes := append(question.Yes, append(question.No, question.NoIdea...)...)
	for _, vote := range votes {
		if vote == userID {
			return true
		}
	}
	return result
}
