package main

type User struct {
    ID          string
    SecretCode  string
    Name        string
    Email       string
    Complaints  []string
}

type Complaint struct {
    ID        string
    Title     string
    Summary   string
    Severity  int
    Resolved  bool
    UserID    string
}