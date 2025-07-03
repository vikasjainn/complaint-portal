// Common/Common_test.go
package Common

import (
	"testing"
)

// TestGenerateID ensures the GenerateID function works as expected.
func TestGenerateID(t *testing.T) {
	id := GenerateID()

	// Test 1: Check that the ID is not empty.
	if id == "" {
		t.Error("Expected GenerateID to return a non-empty string, but it was empty.")
	}

	// Test 2: Check that the ID has the correct length.
	// 4 bytes of random data encoded as hex should result in an 8-character string.
	expectedLength := 8
	if len(id) != expectedLength {
		t.Errorf("Expected ID to have length %d, but got %d", expectedLength, len(id))
	}
}

// TestGenerateSecretCode ensures the GenerateSecretCode function works as expected.
func TestGenerateSecretCode(t *testing.T) {
	secret := GenerateSecretCode()

	// Test 1: Check that the secret code is not empty.
	if secret == "" {
		t.Error("Expected GenerateSecretCode to return a non-empty string, but it was empty.")
	}

	// Test 2: Check that the secret code has the correct length.
	// 6 bytes of random data encoded as hex should result in a 12-character string.
	expectedLength := 12
	if len(secret) != expectedLength {
		t.Errorf("Expected secret code to have length %d, but got %d", expectedLength, len(secret))
	}
}