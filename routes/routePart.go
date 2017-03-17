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
	err = datasource.AddPart(&newPart)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, fileHeaders := range req.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			for _, f := range newPart.Files {
				if f.FileName == fileHeader.Filename {
					file, _ := fileHeader.Open()
					datasource.SavePartFile(companyID, datasource.ObjectIDToID(newPart.ID), f, file)
				}
			}
		}

	}

	w.WriteHeader(http.StatusOK)
})

// ProcessParts - Save the time the part was processed
var ProcessParts = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	partFileID := vars["partFileID"]
	partID := vars["partID"]
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

	err := datasource.ProcessPart(partID, partFileID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not set part as processed"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
})

// GetPartById - Get a Part by PartID
var GetPartById = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	partID := vars["partID"]
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

	part, err := datasource.GETPart(partID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not found"))
		return
	}

	response, err := json.Marshal(&part)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error parsing to json"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

var DeletePart = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	partID := vars["partID"]
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

	err := datasource.DeletePart(partID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not delete Part"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
})

var GetBarcodePart = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	barcode := vars["barcode"]
	machineID := vars["machineID"]
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

	part, err := datasource.GETPartByBarcode(machineID, barcode)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	response, err := json.Marshal(&part)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error returning batches"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
})

var GetPartsByMachine = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	machineID := vars["machineID"]
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

	part, err := datasource.GETPartsByMachine(machineID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	response, err := json.Marshal(&part)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error returning batches"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
})

var GetPartsByCompany = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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

	part, err := datasource.GetPartsByCompany(companyID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	response, err := json.Marshal(&part)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error returning parts"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
})
