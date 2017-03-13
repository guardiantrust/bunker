package datasource

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

const mongoAddress string = "ds046549.mlab.com:46549"
const MongoDatabase string = "bunker-dev"
const mongoUser string = "aaronpeterson3"
const mongoPassword string = "Ethan000111"

// UserCollection - string to store the name of the users collection
const UserCollection string = "Users"

// CompanyCollection - string to store the name of the companies collection
const CompanyCollection string = "Companies"

// BatchCollection - string to store the name of the batches collection
const BatchCollection string = "Batches"

// PartCollection - string to store the name of the parts collection
const PartCollection string = "Parts"

// MachineCollection - string to store the name of the machines collection
const MachineCollection string = "Machines"

// TokenCollection - string to store te name of the token collection
const TokenCollection string = "Tokens"

var mainSession mgo.Session

//GetDBSession - Get a new copied session to the database
func GetDBSession() *mgo.Session {
	return mainSession.Copy()
}

//CloseDBSession - Close the copied Database session
func CloseDBSession(copiedSession *mgo.Session) {
	copiedSession.Close()
}

// ConnectDatabase - Make base connection to database
func ConnectDatabase() {

	//Setup Mongodb connection
	dialerInfo := &mgo.DialInfo{
		Addrs:    []string{mongoAddress},
		Database: MongoDatabase,
		Username: mongoUser,
		Password: mongoPassword,
	}

	//Get Mongodb session
	mongoSession, err := mgo.DialWithInfo(dialerInfo)

	if err != nil {
		fmt.Println("Error creating db session: ")
	}

	//Set session mode
	mongoSession.SetMode(mgo.Monotonic, true)

	mainSession = *mongoSession.Copy()
	defer mongoSession.Close()
}

//DisconnectDatabase - Disconnect from database
func DisconnectDatabase() {
	mainSession.Close()
}
