package datasource

import (
	models "bunker/models"
	"fmt"

	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// StoreToken - Store the new token in the db to validate later
func StoreToken(token *models.Token, session *mgo.Session) {

	col := session.DB(MongoDatabase).C(TokenCollection)
	DeleteUserToken(token.UserID, session)
	col.Insert(&token)
}

// CheckToken - Check the values of the token, against the values saved in the db
func CheckToken(userID string, companyID string, token string) (bool, error) {
	session := GetDBSession()
	defer session.Close()

	col := session.DB(MongoDatabase).C(TokenCollection)
	var savedToken models.Token
	fmt.Println(userID, " ", companyID, " ", token)
	err := col.Find(bson.M{"userID": userID, "companyID": companyID, "token": token}).One(&savedToken)

	if err != nil {
		return false, errors.New("Token not valid")
	}

	if savedToken.ExpiresAt < float64(time.Now().Unix()) {
		DeleteUserToken(userID, session)
		return false, errors.New("Token expired")
	}

	if savedToken.IssuedAt > float64(time.Now().Unix()) {
		DeleteUserToken(userID, session)
		return false, errors.New("Token parameters invalid")
	}

	if len(savedToken.UserID) == 0 {
		DeleteUserToken(userID, session)
		return false, errors.New("UserId invalid")
	}

	if len(savedToken.CompanyID) == 0 {
		DeleteUserToken(userID, session)
		return false, errors.New("ComapnyID invalid")
	}

	if len(savedToken.Token) == 0 {
		DeleteUserToken(userID, session)
		return false, errors.New("Token Invalid")
	}

	return true, nil
}

// DeleteUserToken - delete all tokens associated with a userID
func DeleteUserToken(userID string, session *mgo.Session) {

	col := session.DB(MongoDatabase).C(TokenCollection)
	col.RemoveAll(bson.M{"userID": userID})
}
