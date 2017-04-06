package routes

import (
	"bunker/datasource"
	"bunker/models"
	"bunker/security"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"bytes"

	"github.com/gorilla/mux"
)

var AddParts = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	companyID := vars["companyID"]

	loggedInUser, _ := security.ValidateToken(req)
	if loggedInUser.CompanyID != companyID {
		// check if user is admin
		user, _ := datasource.GetUser(loggedInUser.UserID)

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
		w.Write([]byte(err.Error()))
		return
	}
	response, _ := json.Marshal(newPart)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
})

var UpdatePart = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	companyID := vars["companyID"]
	partID := vars["partID"]

	loggedInUser, _ := security.ValidateToken(req)
	if loggedInUser.CompanyID != companyID {
		// check if user is admin
		user, _ := datasource.GetUser(loggedInUser.UserID)

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
	err = datasource.UpdatePart(partID, &newPart)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
})

// ProcessParts - Save the time the part was processed
var ProcessParts = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	partFileID := vars["fileID"]
	partID := vars["partID"]
	companyID := vars["companyID"]

	loggedInUser, _ := security.ValidateToken(req)
	if loggedInUser.CompanyID != companyID {
		// check if user is admin
		user, _ := datasource.GetUser(loggedInUser.UserID)

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
		user, _ := datasource.GetUser(loggedInUser.UserID)

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
		user, _ := datasource.GetUser(loggedInUser.UserID)

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
		user, _ := datasource.GetUser(loggedInUser.UserID)

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
		user, _ := datasource.GetUser(loggedInUser.UserID)

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
		user, _ := datasource.GetUser(loggedInUser.UserID)

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

var GetFile = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fileID := vars["fileID"]
	companyID := vars["companyID"]

	loggedInUser, _ := security.ValidateToken(req)
	if loggedInUser.CompanyID != companyID {
		// check if user is admin
		user, _ := datasource.GetUser(loggedInUser.UserID)

		if user.Userlevel != models.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
	}

	file, partFile, err := datasource.GetPartFile(fileID)

	if err != nil {
		// file not found
	}

	var contentType = ""
	if len(file) <= 512 {
		contentType = http.DetectContentType(file)
	} else {
		fileHeader := make([]byte, 512)
		for i := range fileHeader {
			fileHeader[i] = file[i]
		}
		contentType = http.DetectContentType(fileHeader)
	}

	fileSize := strconv.FormatInt(int64(binary.Size(file)), 10)
	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+partFile.FileName+"."+partFile.FileExtension)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fileSize)
	sendFile := bytes.NewReader(file)
	//Send the file
	io.Copy(w, sendFile) //'Copy' the file to the client

})

var AddFile = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fileID := vars["fileID"]
	companyID := vars["companyID"]

	loggedInUser, _ := security.ValidateToken(req)
	if loggedInUser.CompanyID != companyID {
		// check if user is admin
		user, _ := datasource.GetUser(loggedInUser.UserID)

		if user.Userlevel != models.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
	}

	for _, fileHeaders := range req.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			datasource.SavePartFile(companyID, fileID, fileHeader.Filename, file)
		}

	}
})
