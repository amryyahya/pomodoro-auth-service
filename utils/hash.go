package main

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// Parameters for Argon2
const (
	saltLength = 16 // Length of the salt in bytes
	keyLength  = 32 // Length of the hash in bytes
)

func GenerateRandomSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

func HashPassword(password string) ([]byte, []byte, error) {
	salt, err := GenerateRandomSalt()
	if err != nil {
		return nil, nil, err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 1, keyLength)
	return hash, salt, nil
}

func VerifyPassword(password string, hash, salt []byte) bool {
	newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 1, keyLength)
	return bytes.Equal(newHash, hash)
}

func main() {
	password := "my_secure_password"

	// Hash the password
	hash, salt, err := HashPassword(password)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		return
	}

	fmt.Printf("Hash: %x\n", hash)
	fmt.Printf("Salt: %x\n", salt)

	// Verify the password
	isValid := VerifyPassword(password, hash, salt)
	if isValid {
		fmt.Println("Password verification successful!")
	} else {
		fmt.Println("Password verification failed!")
	}
}
