package security

import (
	"math/rand"
	"time"

	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var key = []byte("16118C5BFC83B205598690AB86535AA4")
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// EncryptString - Apply site standard encryption to a string
func EncryptString(val []byte) (string, error) {

	hashed, _ := bcrypt.GenerateFromPassword(val, 10)
	return string(hashed), nil
}

// DecryptString - Decrypet a string with site standard encryption
func DecryptString(hashedPassword []byte, val []byte) bool {

	err := bcrypt.CompareHashAndPassword(hashedPassword, val)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Password Success")
	return true
}

// RandStringBytesMaskImprSrc - Create a random string
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
