package utils

// This will handle all aspects of authenticating users in our system
// For password managing/salting I used:
// http://austingwalters.com/building-a-web-server-in-go-salting-passwords/

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

const (
	SaltLength = 64
	// On a scale of 3 - 31, how intense Bcrypt should be
	EncryptCost = 14
)

// This is returned when a new hash + salt combo is generated
// no need for new struct
// type Password struct {
//	hash string
//	salt string
//}

// this handles taking a raw user password and making in into something safe for
// storing in our DB
func hashPassword(salted_pass string) string {

	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(salted_pass), EncryptCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashed_pass)
}

// Handles merging together the salt and the password
func Combine(salt string, raw_pass string) string {

	// concat salt + password
	pieces := []string{salt, raw_pass}
	salted_password := strings.Join(pieces, "")
	return salted_password
}

// Generates a random salt using DevNull
func generateSalt() string {

	// Read in data
	data := make([]byte, SaltLength)
	_, err := rand.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	// Convert to a string
	salt := string(data[:])
	return salt
}

// Handles create a new hash/salt combo from a raw password as inputted
// by the user
//func CreatePassword(raw_pass string) *Password {

func GenerateHashAndSalt(raw_pass string) (string, string) {
	//password := new(Password)
	//password.salt = generateSalt()
	salt := generateSalt()
	salted_pass := Combine(salt, raw_pass)
	//password.hash = hashPassword(salted_pass)
	password_hash := hashPassword(salted_pass)
	return password_hash, salt
	//return password
}
