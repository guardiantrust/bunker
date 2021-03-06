package routes

import (
	"fmt"
	"log"
	"net/http"

	"bunker/datasource"
	security "bunker/security"

	"github.com/gorilla/handlers"
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
	router.Handle("/api/v1/companies/{companyID}", AuthorizedHandler(GETCompany))
	router.Handle("/api/v1/companies/{companyID}", AuthorizedHandler(PUTCompany)).Methods("PUT")
	router.Handle("/api/v1/companies/", GETAllCompany).Methods("GET", "OPTIONS")
	//router.Handle("/api/v1/companies/", AuthorizedHandler(POSTCompany)).Methods("POST")
	router.Handle("/api/companies/{companyID}/users", AuthorizedHandler(GETAllUsers)).Methods("GET")
	router.Handle("/api/v1/companies/{companyID}/parts", AuthorizedHandler(GetPartsByCompany)).Methods("GET")
	router.Handle("/api/v1/companies/{companyID}/parts/{partID}", AuthorizedHandler(GetPartById)).Methods("GET")
	router.Handle("/api/v1/companies/{companyID}/parts/{partID}", AuthorizedHandler(UpdatePart)).Methods("PUT")
	router.Handle("/api/v1/companies/{companyID}/parts", AuthorizedHandler(AddParts)).Methods("POST")
	router.Handle("/api/v1/companies/{companyID}/users", AuthorizedHandler(POSTUser)).Methods("POST")
	router.Handle("/api/v1/companies/{companyID}/batches", AuthorizedHandler(GETCompanyBatches)).Methods("GET")
	router.Handle("/api/v1/companies/{companyId}/batches", AuthorizedHandler(POSTBatch)).Methods("POST")
	router.Handle("/api/v1/companies/{companyID}/batches/{batchId}", AuthorizedHandler(GETBatch)).Methods("GET")
	router.Handle("/api/v1/companies/{companyID}/batches", AuthorizedHandler(PUTBatch)).Methods("PUT")
	router.Handle("/api/v1/companies/{companyID}/batches/{batchID}", AuthorizedHandler(DELETEBatch)).Methods("DELETE")
	router.Handle("/api/v1/companies/{companyID}/machines", AuthorizedHandler(GetMachineByCompany)).Methods("GET")
	router.Handle("/api/v1/companies/{companyID}/machines", AuthorizedHandler(AddMachine)).Methods("POST")
	router.Handle("/api/v1/companies/{companyID}/machines/{machineID}", AuthorizedHandler(GETMachine)).Methods("GET")
	router.Handle("/api/v1/companies/{companyId}/machines/{machineID}", AuthorizedHandler(UpdateMachine)).Methods("PUT")
	router.Handle("/api/v1/companies/{companyID}/machines/{machineID}", AuthorizedHandler(InactivateMachine)).Methods("DELETE")
	router.Handle("/api/v1/companies/{companyID}/machines", AuthorizedHandler(GetMachineByCompany)).Methods("GET")
	router.Handle("/api/v1/companies/{companyID}/machines/{machineID}/parts", AuthorizedHandler(GetPartsByMachine)).Methods("GET")
	router.Handle("/api/v1/companies/{companyID}/machines/{machineID}/parts/{partID}", AuthorizedHandler(GetPartById)).Methods("GET")
	router.Handle("/api/v1/companies/{companyID}/machines/{machineID}/parts/{partId}/process/{fileID}", AuthorizedHandler(ProcessParts)).Methods("POST")
	router.Handle("/api/v1/companies/{companyID}/machines/{machineID}/parts/barcodes/{barcode}", AuthorizedHandler(GetBarcodePart)).Methods("GET")
	router.Handle("/api/v1/companies/{companyID}/machines/{machineID}/parts/files/{fileID}", AuthorizedHandler(GetFile)).Methods("GET")
	router.Handle("/api/v1/companies/{companyID}/machines/{machineID}/parts/files/{fileID}", AuthorizedHandler(AddFile)).Methods("POST")
}

// Listen Listen on port
func Listen() {
	if router != nil {
		//headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "text/plain; charset=utf-8"})
		corsObj := handlers.AllowedOrigins([]string{"*"})
		//methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
		log.Fatal(http.ListenAndServe(":5050", handlers.CORS(corsObj)(router)))
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
