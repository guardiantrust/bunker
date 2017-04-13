package routes

import (
	"bunker/datasource"
	"bunker/models"
	"bunker/security"
	"encoding/json"

	"io"
	"io/ioutil"
	"net/http"

	"fmt"

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
	tempSession := datasource.GetDBSession()
	defer datasource.CloseDBSession(tempSession)
	vars := mux.Vars(req)
	fileID := vars["fileID"]

	fmt.Println("Open File")
	file, err := datasource.GetPartFile(fileID, tempSession)
	fmt.Println("Return File")
	if err != nil {
		// file not found
	}

	//	var contentType = ""
	//if file.Size() <= 512 {
	//	fileHeader := make([]byte, file.Size())
	//	file.Read(fileHeader)
	//	contentType = http.DetectContentType(fileHeader)
	//} else {
	//	fileHeader := make([]byte, 512)
	//	file.Read(fileHeader)
	//	contentType = http.DetectContentType(fileHeader)
	//}

	//Send the file

	//Send the headers

	//w.Header().Set("Content-Disposition", "inline; filename=\""+file.Name()+"\"")
	//w.Header().Set("Content-Type", contentType)
	//w.Header().Set("Content-Length", strconv.FormatInt(file.Size(), 10))
	//if w.Header().Get("Content-Disposition") == "" {
	//	w.Header().Add("Content-Disposition", "inline; filename=\""+file.Name()+"\"")
	//} else {
	//	w.Header().Set("Content-Disposition", "inline; filename=\""+file.Name()+"\"")
	//}

	//if w.Header().Get("Content-Length") == "" {
	//	w.Header().Add("Content-Length", strconv.FormatInt(file.Size(), 10))
	//} else {
	//		w.Header().Set("Content-Length", strconv.FormatInt(file.Size(), 10))
	//	}
	defer file.Close()
	_, err = io.Copy(w, file) //'Copy' the file to the client

	if err != nil {
		fmt.Println(err)
	}

	//w.WriteHeader(http.StatusOK)
	fmt.Println("Done")
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

	fmt.Println("Start to read")
	if err := req.ParseMultipartForm(5 << 20); err != nil {
		fmt.Println(err)
		return
	}

	for _, fileHeaders := range req.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			fmt.Println("Open")
			file, err := fileHeader.Open()
			if err != nil {
				fmt.Println(err)
			}

			err = datasource.SavePartFile(companyID, fileID, fileHeader.Filename, file)

			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Done")
		}

	}
})
