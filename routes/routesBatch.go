package routes

import (
	datasource "bunker/datasource"
	"bunker/models"
	"bunker/security"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// GETBatch - Get a company
var GETBatch = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	companyID := vars["companyID"]
	batchID := vars["batchID"]

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

	batch, err := datasource.GetBatchByID(batchID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Batch not found"))
		return
	}
	response, _ := json.Marshal(&batch)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

// GETCompanyBatches - Return all batches for a company
var GETCompanyBatches = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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

	batches, err := datasource.GetBatchesByCompanyID(companyID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Batch no found"))
		return
	}

	response, err := json.Marshal(&batches)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error returning batches"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
})

// POSTBatch - add a batch
var POSTBatch = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

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

	// Read body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	// copy to model
	var batch models.Batch
	err = json.Unmarshal(body, &batch)

	if err != nil {
		panic(err)
	}

	// add created time
	batch.Created = time.Now()

	err = datasource.AddBatch(&batch)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request - check Batch model format"))
		return
	}

	batch, err = datasource.GetBatchByID(datasource.ObjectIDToID(batch.BatchID))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error retreiving saved Batch"))
		return
	}

	response, _ := json.Marshal(&batch)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
})

// PUTBatch - add a batch
var PUTBatch = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	companyID := vars["companyID"]
	batchID := vars["batchID"]

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

	// Read body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	// copy to model
	var batch models.Batch
	err = json.Unmarshal(body, &batch)

	if err != nil {
		panic(err)
	}

	err = datasource.UpdateBatch(batchID, &batch)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request - check Batch model format"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
})

// DELETEBatch - Delete a batch by ID
var DELETEBatch = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	companyID := vars["companyID"]
	batchID := vars["batchID"]

	// Test access

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

	err := datasource.DeleteBatchByID(batchID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not delete the batch"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
})
