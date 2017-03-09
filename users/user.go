package users

import (
	"time"
)

// SecurityLevel of the user
type SecurityLevel uint8

// SecurityLevels possible
const (
	ADMIN SecurityLevel = iota + 1
	MANAGER
	SUPERVISOR
	OPERATOR
	OFFICE
)

// User structure
type User struct {
	FirstName  string        `json:"UserName"`
	LastName   string        `json:"_id"`
	Email      string        `json:"Email"`
	Department string        `json:"department"`
	Security   SecurityLevel `json:"securityLevel"`
	LastLogin  time.Time     `json:"lastLogin"`
}

// Init User
func (u *User) Init(fName string, lName string, email string, dept string, sec SecurityLevel, ll time.Time) {
	u.FirstName = fName
	u.LastName = lName
	u.Email = email
	u.Department = dept
	u.Security = sec
	u.LastLogin = ll
}

// GetSecurityName - Returns the text of the security level
func (u *User) GetSecurityName() string {
	switch sl := u.Security; sl {
	case 1:
		return "ADMIN"
	case 2:
		return "MANAGER"
	case 3:
		return "SUPERVISOR"
	case 4:
		return "MANAGER"
	case 5:
		return "OPERATOR"
	case 6:
		return "OFFICE"
	default:
		return "NADA"
	}
}
