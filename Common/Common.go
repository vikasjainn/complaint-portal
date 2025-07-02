// Common/Common.go
package Common

import (
	"crypto/rand"
	"encoding/hex"
)

// --- Struct Definitions ---

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

// --- Utility Functions ---

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
