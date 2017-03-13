package security

import (
	"golang.org/x/crypto/bcrypt"
)

var key = []byte("16118C5BFC83B205598690AB86535AA4")

// EncryptString - Apply site standard encryption to a string
func EncryptString(val []byte) (string, error) {

	hashed, _ := bcrypt.GenerateFromPassword(val, 10)
	return string(hashed), nil
}

// DecryptString - Decrypet a string with site standard encryption
func DecryptString(hashedPassword []byte, val []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, val)
	if err != nil {
		return false
	}
	return true
}
