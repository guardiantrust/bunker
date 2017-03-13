package routes

import (
	datasource "bunker/datasource"
	"bunker/models"
	security "bunker/security"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var router *mux.Router

// SetupRoutes - Build all the routes for the service
func SetupRoutes() {
	router = mux.NewRouter()
	router.HandleFunc("/status", StatusHandler).Methods("GET")
	router.HandleFunc("/token", GetTokenHandler).Methods("POST")
	router.HandleFunc("/check-token", CheckTokenHandler)
	router.HandleFunc("/logout", LogoutHandler)
	router.HandleFunc("/api/v1/users/{userId}", NotImplemented).Methods("GET")
	router.HandleFunc("/api/v1/users/{userId}", NotImplemented).Methods("PUT")
	router.HandleFunc("/api/v1/companies/{companyId}", NotImplemented).Methods("GET")
	router.HandleFunc("/api/v1/companies/{companyId}", NotImplemented).Methods("PUT")
	router.HandleFunc("/api/v1/companies/", POSTCompany).Methods("POST")
	router.HandleFunc("/api/companies/{companyId}/users", POSTUser).Methods("GET")
	router.HandleFunc("/api/v1/companies/{companyId}/users", POSTUser).Methods("POST")
	router.HandleFunc("/api/v1/companies/{companyId}/batches", NotImplemented).Methods("GET", "POST")
	router.HandleFunc("/api/v1/companies/{companyId}/batches/{batchId}", NotImplemented).Methods("GET", "DELETE", "PUT")
	router.HandleFunc("/api/v1/companies/{companyId}/machines", NotImplemented).Methods("GET", "POST")
	router.HandleFunc("/api/v1/companies/{companyId}/machines/{machineId}", NotImplemented).Methods("GET", "PUT", "DELETE")

}

// Listen Listen on port
func Listen() {
	if router != nil {
		log.Fatal(http.ListenAndServe(":5050", router))
	} else {
		fmt.Println("Error starting port")
	}
}

// GetTokenHandler - Create a token for an authenticated user
var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	userName := r.Header.Get("username")
	password := r.Header.Get("password")
	session := datasource.GetDBSession()
	defer session.Close()

	validUser, err := datasource.ValidateUser(userName, password, session)

	if err != nil {
		fmt.Println("Error logging in ", userName)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized"))
		return
	}
	userToken := new(models.Token)

	userToken.ExpiresAt = float64(time.Now().Add(time.Hour * 1).Unix())
	userToken.IssuedAt = float64(time.Now().Unix())
	userToken.Token = security.GetToken(&validUser, userToken.ExpiresAt, userToken.IssuedAt)
	userToken.CompanyID = validUser.CompanyID.String()
	userToken.UserID = validUser.ID.String()

	go datasource.UpdateUserLastLogin(&validUser, session)

	go datasource.StoreToken(userToken, session)

	// If we get here, we should have a valid user
	w.Write([]byte(userToken.Token))
})

// CheckTokenHandler - Check the if a token is valid
var CheckTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	validToken, err := security.ValidateToken(strings.Replace(r.Header.Get("Authorization"), "bearer ", "", -1))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	valid, dataErr := datasource.CheckToken(validToken.UserID, validToken.CompanyID, validToken.Token)

	if dataErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(dataErr.Error()))
		return
	}

	if valid {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Failure to validate token"))

})

var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Checking Status")
	w.Write([]byte("All is good!"))

})

var LogoutHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//Logout - kill token
	w.Write([]byte("Logout"))
})

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//Not needed
	w.Write([]byte("Not Implemented"))
})
