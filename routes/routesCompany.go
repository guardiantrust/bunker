package routes

import (
	datasource "bunker/datasource"
	models "bunker/models"
	security "bunker/security"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// GETCompany - Get a company
var GETCompany = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	companyID := vars["companyID"]

	loggedInUser, _ := security.ValidateToken(req)
	if loggedInUser.CompanyID != companyID {
		// check if user is admin
		user, _ := datasource.GetUser(datasource.IDToObjectID(loggedInUser.UserID))

		if user.Userlevel != models.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
	}

	company, err := datasource.GetCompany(companyID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Company not found"))
		return
	}
	response, _ := json.Marshal(&company)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

// POSTCompany - Add a company
var POSTCompany = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// Read body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	// copy to model
	var company models.Company
	err = json.Unmarshal(body, &company)

	if err != nil {
		panic(err)
	}

	// add created time
	company.Created = time.Now()
	// push to db
	err = datasource.AddCompany(&company)
	if err != nil {
		// error will be thrown if name is not available
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Company name already taken"))
		return
	}
	// retrieve new record
	companies := datasource.GetCompanyByName(company.Name)
	// write response
	response, _ := json.Marshal(&companies)
	w.Write(response)
})

var PUTCompany = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// Read body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var company models.Company
	err = json.Unmarshal(body, &company)

	//check valid object
	//userName := req.Header.Get("username")

})

// DeleteCompany - Inactivates a company in the db
var DeleteCompany = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	variables := mux.Vars(req)
	companyID := variables["companyID"]

	err := datasource.InActivateCompany(companyID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}
	w.WriteHeader(http.StatusAccepted)
})
