package routes

import (
	"bunker/datasource"
	"bunker/models"
	"bunker/security"
	"net/http"
	"strings"
	"sync"
	"time"
)

// GetTokenHandler - Create a token for an authenticated user
var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup

	userName := r.Header.Get("username")
	password := r.Header.Get("password")

	validUser, err := datasource.ValidateUser(userName, password)

	if err != nil {
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

	//Update the users last login
	go func() {
		defer wg.Done()
		wg.Add(1)
		datasource.UpdateUserLastLogin(&validUser)
	}()

	// Update the token in the db
	go func() {
		defer wg.Done()
		wg.Add(1)
		datasource.StoreToken(userToken)
	}()

	// If we get here, we should have a valid user
	w.Write([]byte(userToken.Token))
	wg.Wait()
})

// CheckTokenHandler - Check the if a token is valid
var CheckTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	validToken, err := security.ValidateToken(r)

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
