package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/denisbrodbeck/machineid"
)

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Encrypt encrypt data using aes
func Encrypt(key []byte, text string) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	plaintext := []byte(fmt.Sprint(text))
	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plaintext))
	cfb.XORKeyStream(cipherText, plaintext)
	return encodeBase64(cipherText)
}

// Decrypt decrypt data using aes
func Decrypt(key []byte, text string) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	cipherText := decodeBase64(text)
	cfb := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(cipherText))
	cfb.XORKeyStream(plaintext, cipherText)
	return string(plaintext)
}

// SumHash generate the hash sum of a text
func SumHash(text string) []byte {
	hash := sha256.New()
	return hash.Sum([]byte(text))
}

// EncodeHash encode the hash sum
func EncodeHash(hash []byte) string {
	return hex.EncodeToString(hash[:])
}

// SumHashMachine generate a user machine ID, write the id and the text using sha256, and return the sum of the hash
func SumHashMachine(text string) ([]byte, error) {
	id, err := machineid.ID()
	if err != nil {
		return nil, err
	}

	hash := sha256.New()
	_, _ = hash.Write([]byte(text))
	_, _ = hash.Write([]byte(id))
	sum := hash.Sum(nil)

	return sum, nil
}
