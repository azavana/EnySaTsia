package routes

import (
	db "voting/database"
	"voting/models"

	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
)

// NewSession --create a session !!! should not create empty session
func NewSession(ctx iris.Context) {
	//define session model
	s := models.Session{}
	//get the data from ctx
	err := ctx.ReadJSON(&s)
	//set state of session to created
	s.State = "created"
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Error while getting data from client")
	} else {
		session := db.Connect()
		defer session.Close()
		c := session.DB(db.Database).C("sessions")
		e := c.Insert(s)
		if e != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON("failed to create session")
		} else {
			ctx.StatusCode(iris.StatusCreated)
			ctx.JSON("Session created")
		}
	}
}

// UpdateSession -- Session name and number of voters can be updated, every change schould be tracked
func UpdateSession(ctx iris.Context) {
	// get session Id and the updated value
	s := models.Session{}
	err := ctx.ReadJSON(&s)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Error while getting data from client")
	}

	// Only change state when session is modifiable
	if IsSessionModifiable(s.ID) {
		session := db.Connect()
		defer session.Close()
		c := session.DB(db.Database).C("sessions")

		// Find Session in the db that has to be changes
		sess := models.Session{}
		c.Find(bson.M{"_id": s.ID}).One(&sess)

		// Performe the change and Update the Session
		sess.Name = s.Name
		e := c.Update(bson.M{"_id": s.ID}, &sess)
		if e == nil {
			ctx.StatusCode(iris.StatusAccepted)
			ctx.JSON("Session updated")
		}
	} else {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON("you can not modify a session that is already opened")
	}

}

//StartSession -- start the session
func StartSession(ctx iris.Context) {
	//get the session id
	s := models.Session{}
	err := ctx.ReadJSON(&s)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Error while getting data from client")
	}

	//Update the db
	session := db.Connect()
	defer session.Close()
	c := session.DB(db.Database).C("sessions")

	// Find Session in the db that has to be changes
	sess := models.Session{}
	c.Find(bson.M{"_id": s.ID}).One(&sess)
	//set state to open
	sess.State = "open"
	e := c.Update(bson.M{"_id": s.ID}, &sess)

	if e == nil {
		ctx.StatusCode(iris.StatusAccepted)
		ctx.JSON("Session is started")
	} else {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Something went wrong")
	}
	//send notification to all participants that session is opened
}

//CloseSession -- archive the actual session information
func CloseSession(ctx iris.Context) {
	//get the session id
	s := models.Session{}
	err := ctx.ReadJSON(&s)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Error while getting data from client")
	}
	//set state to open
	s.State = "close"

	//Update the db
	session := db.Connect()
	defer session.Close()
	c := session.DB(db.Database).C("sessions")

	// Find Session in the db that has to be changes
	sess := models.Session{}
	c.Find(bson.M{"_id": s.ID}).One(&sess)
	//set state to open
	sess.State = "close"
	e := c.Update(bson.M{"_id": s.ID}, &sess)

	if e == nil {
		ctx.StatusCode(iris.StatusAccepted)
		ctx.JSON("Session is closed")
	} else {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Something went wrong")
	}
	//send notification to all participant that session is closed
}

//GetAllSession -- get all session in db
func GetAllSession(ctx iris.Context) {
	var ss = []models.Session{}
	session := db.Connect()
	defer session.Close()
	c := session.DB(db.Database).C("sessions")
	e := c.Find(bson.M{}).All(&ss)

	if e == nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ss)
	} else {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Something went wrong")
	}
}

//GetSession --get a specific session
func GetSession(ctx iris.Context) {
	s := models.Session{}
	session := db.Connect()
	defer session.Close()
	c := session.DB(db.Database).C("sessions")
	e := c.FindId(bson.ObjectIdHex(ctx.Params().Get("id"))).One(&s)

	if e == nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(s)
	} else {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON("Something went wrong")
	}
}
