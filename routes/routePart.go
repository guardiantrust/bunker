package routes

import (
	"bunker/datasource"
	"bunker/models"
	"bunker/security"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

var AddParts = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

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

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var newPart models.Part
	err = json.Unmarshal(body, &newPart)

	if err != nil {
		panic(err)
	}

	// Save Part
	for _, fileHeaders := range req.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			for _, f := range newPart.Files {
				if f.FileName == fileHeader.Filename {
					file, _ := fileHeader.Open()

				}
			}
		}

	}

})
