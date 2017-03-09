package datasource

import (
	models "bunker/models"
	security "bunker/security"
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ValidateUser - Validate a user password
func ValidateUser(userName string, password string) (models.User, error) {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)

	user := models.User{}
	coll := tempSession.DB(MongoDatabase).C(UserCollection)
	err := coll.Find(bson.M{"username": userName}).One(&user)

	if err != nil || len(user.Username) == 0 {
		fmt.Println("User Not found")
		return models.User{}, errors.New("User Not Found")
	}
	valid := security.DecryptString([]byte(user.Password), []byte(password))
	if valid {
		fmt.Println("Valid Password")
		return user, nil
	}
	return user, errors.New("Invalid Password")
}

// usernameAvailable - Check if a username is not currently used true = available, false = already taken
func usernameAvailable(userName string) bool {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)

	coll := tempSession.DB(MongoDatabase).C(UserCollection)
	totalUsed, _ := coll.Find(bson.M{"username": userName}).Count()

	return totalUsed == 0
}

// AddUser - Add a new user
func AddUser(user *models.User) error {
	if !usernameAvailable(user.Username) {
		return errors.New("Username already taken")
	}

	encryptedPassword, _ := security.EncryptString([]byte(user.Password))

	user.Password = encryptedPassword
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(UserCollection)
	coll.Insert(&user)
	fmt.Println(user)
	return nil
}

// UpdateUser - Update an existing user
func UpdateUser(user *models.User, session *mgo.Session) {

}

// GetUser - Return a user by userId
func GetUser(userID string, session *mgo.Session) models.User {
	returnUser := new(models.User)

	return *returnUser

}

//GetUserByUserName - return a user by username
func GetUserByUserName(name string) []models.User {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(UserCollection)
	var users []models.User
	coll.Find(bson.M{"username": name}).All(&users)

	return users
}

// GetAllUsers - Return all users belonging to a companyID
func GetAllUsers(companyID string) []models.User {
	tempSession := GetDBSession()
	defer CloseDBSession(tempSession)
	coll := tempSession.DB(MongoDatabase).C(UserCollection)
	users := []models.User{}
	coll.Find(bson.M{"companyId": companyID}).All(&users)

	return users
}

// DeleteUser - Delete (inactivate) a user by userID
func DeleteUser(userID string, session *mgo.Session) {

}
