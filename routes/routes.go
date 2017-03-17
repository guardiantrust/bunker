package routes

import (
	"fmt"
	"log"
	"net/http"

	"bunker/datasource"
	security "bunker/security"

	"github.com/gorilla/mux"
)

var router *mux.Router

// SetupRoutes - Build all the routes for the service
func SetupRoutes() {
	router = mux.NewRouter()
	router.HandleFunc("/v1/status", StatusHandler).Methods("GET")
	router.HandleFunc("/v1/token", GetTokenHandler).Methods("POST")
	router.HandleFunc("/v1/check-token", CheckTokenHandler)
	router.Handle("/v1/logout", AuthorizedHandler(LogoutHandler))
	router.Handle("/api/v1/users/{userID}", AuthorizedHandler(GETUser)).Methods("GET")
	router.Handle("/api/v1/users/{userID}", AuthorizedHandler(PUTUser)).Methods("PUT")
	router.Handle("/api/v1/companies/{companyID}", AuthorizedHandler(GETCompany)).Methods("GET")
	router.Handle("/api/v1/companies/{companyID}", AuthorizedHandler(GETCompany)).Methods("PUT")
	router.Handle("/api/v1/companies/", AuthorizedHandler(POSTCompany)).Methods("POST")
	router.Handle("/api/companies/{companyID}/users", AuthorizedHandler(GETAllUsers)).Methods("GET")
	router.Handle("/api/v1/companies/{companyID}/users", AuthorizedHandler(POSTUser)).Methods("POST")
	router.Handle("/api/v1/companies/{companyID}/batches", NotImplemented).Methods("GET", "POST")
	router.Handle("/api/v1/companies/{companyID}/batches/{batchId}", NotImplemented).Methods("GET", "DELETE", "PUT")
	router.Handle("/api/v1/companies/{companyID}/machines", NotImplemented).Methods("GET", "POST")
	router.Handle("/api/v1/companies/{companyID}/machines/{machineId}", NotImplemented).Methods("GET", "PUT", "DELETE")

}

// Listen Listen on port
func Listen() {
	if router != nil {
		log.Fatal(http.ListenAndServe(":5050", router))
	} else {
		fmt.Println("Error starting port")
	}
}

// AuthorizedHandler - secure endpoint with token authentication
func AuthorizedHandler(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validToken, err := security.ValidateToken(r)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unauthorized"))
			return
		}

		valid, dataErr := datasource.CheckToken(validToken.UserID, validToken.CompanyID, validToken.Token)

		if dataErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(dataErr.Error()))
			return
		}

		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("UnAuthorized"))
			return
		}

		handler.ServeHTTP(w, r)
	})

}

// StatusHandler - check the status of the service
var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Checking Status")
	w.Write([]byte("All is good!"))

})

// LogoutHandler - Logout the current user
var LogoutHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//Logout - kill token
	w.Write([]byte("Logout"))
})

// NotImplemented - used for a place holder for new services being created
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//Not needed
	w.Write([]byte("Not Implemented"))
})
