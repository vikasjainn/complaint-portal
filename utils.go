package main

import (
    "crypto/rand"
    "encoding/hex"
)

func generateID() string {
    b := make([]byte, 4)
    rand.Read(b)
    return hex.EncodeToString(b)
}

func generateSecretCode() string {
    b := make([]byte, 6)
    rand.Read(b)
    return hex.EncodeToString(b)
}