package security

import (
	"golang.org/x/crypto/bcrypt"
)

var key = []byte("16118C5BFC83B205598690AB86535AA4")

// EncryptString - Apply site standard encryption to a string
func EncryptString(val []byte) (string, error) {

	hashed, _ := bcrypt.GenerateFromPassword(val, 10)
	return string(hashed), nil

	/*block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	b := base64.StdEncoding.EncodeToString(val)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return string(ciphertext), nil
	*/
}

// DecryptString - Decrypet a string with site standard encryption
func DecryptString(hashedPassword []byte, val []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, val)
	if err != nil {
		return false
	}
	return true
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.

	/*
		block, err := aes.NewCipher(key)
		if err != nil {
			return "", err
		/}
		if len(val) < aes.BlockSize {
			return "", errors.New("ciphertext too short")
		}
		iv := val[:aes.BlockSize]
		val = val[aes.BlockSize:]
		cfb := cipher.NewCFBDecrypter(block, iv)
		cfb.XORKeyStream(val, val)
		data, err := base64.StdEncoding.DecodeString(string(val))
		if err != nil {
			return "", err
		}
		return string(data), ni*/
}
