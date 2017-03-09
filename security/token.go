package security

import (
	"errors"
	"strconv"
	"time"

	models "bunker/models"

	"fmt"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

var ourLittleSecret = []byte("d0uble-d@re")

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return ourLittleSecret, nil
	},
	SigningMethod: jwt.SigningMethodHS512,
})

// GetToken - Get a token
func GetToken(user *models.User) string {
	// Check userName and password

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["iat"] = time.Now().Unix()
	claims["username"] = user.Username
	claims["companyID"] = user.CompanyID.String()
	claims["roles"] = user.Roles

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	createdToken, _ := token.SignedString([]byte(ourLittleSecret))
	//	fmt.Println("Token:" + createdToken)
	return createdToken
}

func ValidateToken(tokenString string) (bool, error) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return ourLittleSecret, nil
	})

	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, errors.New("token invalid")
	}

	err = token.Claims.Valid()

	if err != nil {
		return false, errors.New("token time invalid")
	}

	expiresAt := token.Claims.(jwt.MapClaims)["exp"].(float64)
	issuedAt := token.Claims.(jwt.MapClaims)["iat"].(float64)
	companyID := token.Claims.(jwt.MapClaims)["companyID"].(string)
	username := token.Claims.(jwt.MapClaims)["username"].(string)

	fmt.Println("Issued At :" + strconv.FormatFloat(issuedAt, 'E', -1, 64))
	fmt.Println("Expires At :" + strconv.FormatFloat(expiresAt, 'E', -1, 64))
	fmt.Println("Username :" + username)
	fmt.Println("CompanyID :" + companyID)

	return true, nil
}
