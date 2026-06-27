package mycrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/argon2"
)

type ArgonParameters struct {
	Memory      uint32
	Iteration   uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func EncryptFile(filePath string, password string, argonParameters ArgonParameters) error {
	key, salt, err := GenerateHash(password, argonParameters)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	ciphertext, err := EncrytWithGCM(data, key)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	output := append(salt, ciphertext...)
	err = os.WriteFile(filePath, output, 0644)
	if err != nil {
		fmt.Println("Error saving file: ", err)
		return err
	}

	return nil
}

func EncrytWithGCM(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func DecryptFile(filePath string, password string, argonParameters ArgonParameters) error {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	salt := data[:argonParameters.SaltLength]
	remaining := data[argonParameters.SaltLength:]

	key := argon2.IDKey([]byte(password), salt, argonParameters.Iteration, argonParameters.Memory, argonParameters.Parallelism, argonParameters.KeyLength)

	plaintext, err := DecryptWithGCM(remaining, key)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, plaintext, 0644)

}

func DecryptWithGCM(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

func GenerateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func GenerateHash(password string, parameter ArgonParameters) (hash []byte, salt []byte, err error) {
	salt, err = GenerateSalt(parameter.SaltLength)
	if err != nil {
		return nil, nil, err
	}

	hash = argon2.IDKey([]byte(password), salt, parameter.Iteration, parameter.Memory, parameter.Parallelism, parameter.KeyLength)
	return hash, salt, nil
}
