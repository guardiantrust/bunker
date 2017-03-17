package routes

import (
	datasource "bunker/datasource"
	models "bunker/models"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

// GETUser - Return a User by UserId
var GETUser = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := bson.ObjectIdHex(vars["userID"])
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	user, err := datasource.GetUser(userID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found with that Id: " + userID.Hex()))
		return
	}

	response, _ := json.Marshal(&user)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

// GETAllUsers - return all the users for a company
var GETAllUsers = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	companyID := vars["companyID"]

	users := datasource.GetAllUsers(companyID)

	response, _ := json.Marshal(&users)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

// POSTUser - Create a new user
var POSTUser = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// Read body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	// copy to model
	var user models.Login
	err = json.Unmarshal(body, &user)
	user.CompanyID = bson.ObjectIdHex(vars["companyID"])

	if err != nil {
		panic(err)
	}

	if !validLoginModel(&user) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrecet Login Model"))
		return
	}

	// push to db
	err = datasource.AddUser(&user)

	if err != nil {
		// error will be thrown if name is not available
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User name already taken"))
		return
	}
	// retrieve new record
	users := datasource.GetUserByUserName(user.Username)

	// write response
	w.Write([]byte(users[0].ID.Hex()))
})

// PUTUser = Update a user in the database
var PUTUser = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["userID"]

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	// copy to model
	var user models.User
	err = json.Unmarshal(body, &user)
	user.ID = datasource.IDToObjectID(userID)
	err = datasource.UpdateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Update did not occur"))
		return
	}
	w.WriteHeader(http.StatusNoContent)

})

// PasswordChange - change the password of a user
var PasswordChange = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["userID"]
	password := vars["password"]

	err := datasource.UpdatePassword(userID, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Password not changed"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
})

func validLoginModel(user *models.Login) bool {

	if !validUserModel(&user.User) {
		return false
	}

	if user.Password == "" {
		return false
	}

	return true
}

func validUserModel(user *models.User) bool {

	if user.Username == "" {
		return false
	}

	if user.FirstName == "" {
		return false
	}

	if user.LastName == "" {
		return false
	}

	if user.Email == "" {
		return false
	}

	return true
}
