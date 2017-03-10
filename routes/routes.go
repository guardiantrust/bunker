package routes

import (
	datasource "bunker/datasource"
	security "bunker/security"
	"fmt"
	"log"
	"net/http"

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
	router.HandleFunc("/api/user/{userId}", NotImplemented).Methods("GET")
	router.HandleFunc("/api/user/{userId}", NotImplemented).Methods("PUT")
	router.HandleFunc("/api/company/{companyId}", NotImplemented).Methods("GET")
	router.HandleFunc("/api/company/{companyId}", NotImplemented).Methods("PUT")
	router.HandleFunc("/api/company", POSTCompany).Methods("POST")
	router.HandleFunc("/api/company/{companyId}/users", POSTUser).Methods("GET")
	router.HandleFunc("/api/company/{companyId}/users", POSTUser).Methods("POST")
	router.HandleFunc("/api/company/{companyId}/batches", NotImplemented).Methods("GET", "POST")
	router.HandleFunc("/api/company/{companyId}/batches/{batchId}", NotImplemented).Methods("GET", "DELETE", "PUT")
	router.HandleFunc("/api/company/{companyId}/machines", NotImplemented).Methods("GET", "POST")
	router.HandleFunc("/api/company/{companyId}/machines/{machineId}", NotImplemented).Methods("GET", "PUT", "DELETE")

}

// Listen Listen on port
func Listen() {
	if router != nil {
		log.Fatal(http.ListenAndServe(":5050", router))
	} else {
		fmt.Println("Error starting port")
	}
}

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	userName := r.Header.Get("username")
	password := r.Header.Get("password")
	validUser, err := datasource.ValidateUser(userName, password)

	if err != nil {
		fmt.Println("Error logging in ", userName)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized"))
		return
	}

	w.Write([]byte(security.GetToken(&validUser)))
})

var CheckTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	fmt.Println(token)
	security.ValidateToken(token)
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
