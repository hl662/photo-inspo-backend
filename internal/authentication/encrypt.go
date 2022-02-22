package authentication

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

func EncryptAES(plaintext string) string {
	// create cipher
	key := os.Getenv("encryptKey")
	c, err := aes.NewCipher([]byte(key))
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return string(ciphertext[:])
}

func DecryptAES(ct string) string {
	data := []byte(ct)
	key := os.Getenv("encryptKey")
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext[:])
}
