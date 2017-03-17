package security

import (
	"errors"
	"net/http"
	"strings"
	"time"

	models "bunker/models"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

var ourLittleSecret = []byte("trIplE-d0g-d@re-y0U")

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return ourLittleSecret, nil
	},
	SigningMethod: jwt.SigningMethodHS512,
})

// GetToken - Get a token
func GetToken(user *models.User, expiresAt float64, issuedAt float64) string {
	// Check userName and password

	claims := make(jwt.MapClaims)
	claims["exp"] = expiresAt
	claims["iat"] = issuedAt
	claims["userID"] = user.ID.String()
	claims["companyID"] = user.CompanyID.String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	createdToken, _ := token.SignedString([]byte(ourLittleSecret))

	return createdToken
}

// ValidateToken - Check the validity of a token
func ValidateToken(request *http.Request) (models.Token, error) {
	tokenString := strings.Replace(request.Header.Get("Authorization"), "bearer ", "", -1)
	token := new(models.Token)
	tok, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return ourLittleSecret, nil
	})

	if err != nil {
		return *token, err
	}

	if !tok.Valid {
		return *token, errors.New("token invalid")
	}

	err = tok.Claims.Valid()

	if err != nil {
		return *token, errors.New("token time invalid")
	}

	token.ExpiresAt = tok.Claims.(jwt.MapClaims)["exp"].(float64)
	token.IssuedAt = tok.Claims.(jwt.MapClaims)["iat"].(float64)
	token.CompanyID = tok.Claims.(jwt.MapClaims)["companyID"].(string)
	token.UserID = tok.Claims.(jwt.MapClaims)["userID"].(string)
	token.Token = tokenString

	if currentTime := int(time.Now().Unix()); int(token.ExpiresAt) < currentTime {
		return *token, errors.New("Expired token")
	} else if int(token.IssuedAt) > currentTime {
		return *token, errors.New("Invalid issued time")
	} else if len(token.CompanyID) <= 10 {
		return *token, errors.New("Missing info")
	} else if len(token.UserID) <= 0 {
		return *token, errors.New("Missing info")
	}

	return *token, nil
}
