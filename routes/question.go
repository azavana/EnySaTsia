package routes

import (
	"fmt"
	db "voting/database"
	"voting/models"

	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
)

//NewQuestion -- new Question
func NewQuestion(ctx iris.Context) {
	//define question model
	question := models.Question{}
	//Get Question from JSON data
	err := ctx.ReadJSON(&question)

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Error while getting data from client")
	}
	//assume that voted is false by creation
	question.Voted = false

	//assume that state is created by creation
	question.State = "created"

	//assume that Yes, No, NoIdea array are empty by creation
	question.Yes = []string{}
	question.No = []string{}
	question.NoIdea = []string{}

	//assume that Session exist and not closed
	sessionExistsAndNotClosed := SessionExistsAndNotClosed(question.Session)

	if sessionExistsAndNotClosed {
		//add question into database
		session := db.Connect()
		defer session.Close()
		c := session.DB(db.Database).C("questions")
		e := c.Insert(question)
		if e != nil {
			fmt.Print(e)
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON("Error while inserting question")
		} else {
			ctx.StatusCode(iris.StatusCreated)
			ctx.JSON("Question created")
		}
	} else {
		ctx.StatusCode(iris.StatusNotAcceptable)
		ctx.JSON("Session might not exist or already closed")
	}
}

//UpdateQuestion -- only Question sentence can be updated and if vote was not yet opened (voted is false)
func UpdateQuestion(ctx iris.Context) {
	//Get question Id and the updated value
	question := models.UpdateQuestion{}
	err := ctx.ReadJSON(&question)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Error while getting data from client")
	}
	//make sure that question is not already voted
	if IsQuestionModifiable(question.ID) {
		session := db.Connect()
		defer session.Close()
		c := session.DB(db.Database).C("questions")

		// Find Question in the db that has to be changed
		qst := models.Question{}
		c.Find(bson.M{"_id": question.ID}).One(&qst)
		//change the question
		qst.Sentence = question.Sentence

		e := c.Update(bson.M{"_id": question.ID}, &qst)
		if e == nil {
			ctx.StatusCode(iris.StatusAccepted)
			ctx.JSON("Question updated")
		}
	} else {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON("Question cannot be updated, it was already voted")
	}
}

//VoteStart -- start vote for a question
func VoteStart(ctx iris.Context) {
	//Get question Id and the updated value
	question := models.Question{}
	err := ctx.ReadJSON(&question)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Error while getting data from client")
	}
	session := db.Connect()
	defer session.Close()
	c := session.DB(db.Database).C("questions")

	// Find Question in the db that has to be started
	qst := models.Question{}
	c.Find(bson.M{"_id": question.ID}).One(&qst)
	if CanOpen(qst) {
		qst.State = "open"
		e := c.Update(bson.M{"_id": question.ID}, &qst)
		if e != nil {
			ctx.StatusCode(iris.StatusConflict)
			ctx.JSON("Something wrong")
		} else {
			ctx.StatusCode(iris.StatusAccepted)
			ctx.JSON("Vote is opened")
		}
	} else {
		ctx.StatusCode(iris.StatusConflict)
		ctx.JSON("Vote cannot be started")
	}
}

//VoteClose -- close vote for a question
func VoteClose(ctx iris.Context) {
	//Get question Id and the updated value
	question := models.Question{}
	err := ctx.ReadJSON(&question)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Error while getting data from client")
	}
	session := db.Connect()
	defer session.Close()
	c := session.DB(db.Database).C("questions")
	// Find Question in the db that has to be closed
	qst := models.Question{}
	c.Find(bson.M{"_id": question.ID}).One(&qst)
	if CanClose(qst) {
		qst.State = "close"
		e := c.Update(bson.M{"_id": question.ID}, &qst)
		if e != nil {
			ctx.StatusCode(iris.StatusConflict)
			ctx.JSON("Something wrong")
		} else {
			ctx.StatusCode(iris.StatusAccepted)
			ctx.JSON("Vote is closed")
		}
	} else {
		ctx.StatusCode(iris.StatusConflict)
		ctx.JSON("Vote cannot be closed")
	}
}

//Vote -- respond the question, the id of the user is added to the yes or no slice
func Vote(ctx iris.Context) {
	//get question Id, voter Id and vote
	vote := models.Vote{}
	err := ctx.ReadJSON(&vote)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Error while getting data from client")
	}
	//search question in db
	session := db.Connect()
	defer session.Close()
	c := session.DB(db.Database).C("questions")
	q := models.Question{}

	//find the session with sessionId
	e := c.Find(bson.M{"_id": vote.Question}).One(&q)

	//add decision in each respective slice
	if IsVoteOpen(q) {
		if !HasVoted(q, vote.User) {
			if e == nil {
				if vote.Decision == "yes" {
					q.Yes = append(q.Yes, vote.User)
				}
				if vote.Decision == "no" {
					q.No = append(q.No, vote.User)
				}
				if vote.Decision == "noIdea" {
					q.NoIdea = append(q.NoIdea, vote.User)
				}
				save := c.Update(bson.M{"_id": q.ID}, &q)
				if save == nil {
					ctx.StatusCode(iris.StatusAccepted)
					ctx.JSON("success")
				} else {
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.JSON("Error while saving data in database")
				}
			} else {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.JSON("Error while getting data from database")
			}
		} else {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON("You already voted !!!")
		}
	} else {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON("Vote is not yet opened or already closed")
	}
}

//GetQuestion -- get the result of the vote of a question. It actually just returns the question with the vote data
func GetQuestion(ctx iris.Context) {
	// get the question id, session id
	questionID := ctx.Params().Get("question")

	question := models.Question{}
	// get session and question from database
	session := db.Connect()
	defer session.Close()
	c := session.DB(db.Database).C("questions")
	e := c.Find(bson.M{"_id": bson.ObjectIdHex(questionID)}).One(&question)
	if e != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Error while getting data from database")
	} else {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(question)
	}
}

// GetQuestionOfSession -- get all questions of a session
func GetQuestionOfSession(ctx iris.Context) {
	sessionID := ctx.Params().Get("session")
	var questions = []models.Question{}
	session := db.Connect()
	defer session.Close()
	c := session.DB(db.Database).C("questions")
	e := c.Find(bson.M{"session": bson.ObjectIdHex(sessionID)}).All(&questions)

	if e == nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(questions)
	} else {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Something went wrong")
	}
}
