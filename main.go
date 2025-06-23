package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    fmt.Println("Starting Complaint Portal API...")

    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/submitComplaint", submitComplaintHandler)
    http.HandleFunc("/getAllComplaintsForUser", getUserComplaintsHandler)
    http.HandleFunc("/getAllComplaintsForAdmin", getAdminComplaintsHandler)
    http.HandleFunc("/viewComplaint", viewComplaintHandler)
    http.HandleFunc("/resolveComplaint", resolveComplaintHandler)

    log.Fatal(http.ListenAndServe(":8080", nil))
}