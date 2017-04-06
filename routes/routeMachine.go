package routes

import (
	datasource "bunker/datasource"
	"bunker/models"
	"bunker/security"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

var GETMachine = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	companyID := vars["companyID"]
	machineID := vars["machineID"]

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

	machine, err := datasource.GetMachineById(companyID, machineID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Could not find machine"))
		return
	}

	response, _ := json.Marshal(&machine)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

var GetMachineByCompany = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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

	machine, err := datasource.GetMachinesByCompany(companyID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Could not find machine"))
		return
	}

	response, _ := json.Marshal(&machine)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

var AddMachine = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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

	// copy to model
	var machine models.Machine
	body, err := ioutil.ReadAll(req.Body)

	err = json.Unmarshal(body, &machine)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong!"))
		return
	}

	err = datasource.AddMachine(&machine)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong saving machine!"))
		return
	}

	response, _ := json.Marshal(machine)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))

})

var InactivateMachine = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	companyID := vars["companyID"]
	machineID := vars["machineID"]

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

	err := datasource.InactivateMachine(companyID, machineID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not inactivate machine"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

})

var UpdateMachine = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	companyID := vars["companyID"]
	machineID := vars["machineID"]

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

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

	// copy to model
	var machine models.Machine
	err = json.Unmarshal(body, &machine)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong!"))
		return
	}

	err = datasource.UpdateMachine(machineID, machine)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong saving machine!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
})
