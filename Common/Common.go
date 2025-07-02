package Common

import (
	"crypto/rand"
	"encoding/hex"
)


type User struct {
	ID         string
	SecretCode string
	Name       string
	Email      string
	Complaints []string
}

type Complaint struct {
	ID       string
	Title    string
	Summary  string
	Severity int
	Resolved bool
	UserID   string
}


func GenerateID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateSecretCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}
