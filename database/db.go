package database

import (
	"log"

	"gopkg.in/mgo.v2"
)

//var urlDb = "localhost:27017"
var urlDb = "mongodb://mongo:27017"

//Database name
var Database = "voting-app"
var mgoSession *mgo.Session

//Connect --connect to database
func Connect() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(urlDb)
		if err != nil {
			log.Fatal("Failed to start mongo db")
		}
	}
	return mgoSession.Clone()
}
