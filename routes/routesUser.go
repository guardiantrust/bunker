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

	if user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Must include a Password"))
		return
	}

	if user.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Must include a Username"))
		return
	}

	if user.FirstName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Must include a First name"))
		return
	}

	if user.LastName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Must include a Last name"))
		return
	}

	if user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Must include an email"))
		return
	}

	if err != nil {
		panic(err)
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
