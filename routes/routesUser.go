package routes

import (
	datasource "bunker/datasource"
	models "bunker/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"fmt"

	"github.com/gorilla/mux"
)

var POSTUser = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Entering POST User")
	vars := mux.Vars(req)
	// Read body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Body Read")
	// copy to model
	var user models.User
	err = json.Unmarshal(body, &user)
	user.CompanyID = bson.ObjectIdHex(vars["companyId"])
	fmt.Println("User Set :" + user.CompanyID)
	if err != nil {
		panic(err)
	}

	// add created time
	user.Created = time.Now()
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
	response, _ := json.Marshal(&users)
	w.Write(response)
})
