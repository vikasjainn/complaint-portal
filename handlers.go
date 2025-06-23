package main

import (
    "encoding/json"
    "net/http"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
        return
    }

    var input struct {
        Name  string
        Email string
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    if input.Name == "" || input.Email == "" {
        http.Error(w, "Name and Email are required", http.StatusBadRequest)
        return
    }

    mu.Lock()
    for _, u := range users {
        if u.Email == input.Email {
            mu.Unlock()
            http.Error(w, "Email already registered", http.StatusConflict)
            return
        }
    }

    user := User{
        ID:         generateID(),
        SecretCode: generateSecretCode(),
        Name:       input.Name,
        Email:      input.Email,
    }

    users[user.SecretCode] = user
    mu.Unlock()

    json.NewEncoder(w).Encode(user)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    secret := r.URL.Query().Get("secretCode")

    mu.Lock()
    user, ok := users[secret]
    mu.Unlock()

    if !ok {
        http.Error(w, "Invalid secret code", http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func submitComplaintHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
        return
    }

    secret := r.URL.Query().Get("secretCode")
    var input struct {
        Title   string
        Summary string
        Severity int
    }
    json.NewDecoder(r.Body).Decode(&input)

    mu.Lock()
    user, ok := users[secret]
    if !ok {
        mu.Unlock()
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    complaint := Complaint{
        ID:       generateID(),
        Title:    input.Title,
        Summary:  input.Summary,
        Severity: input.Severity,
        UserID:   user.ID,
    }

    complaints[complaint.ID] = complaint
    user.Complaints = append(user.Complaints, complaint.ID)
    users[secret] = user
    mu.Unlock()

    json.NewEncoder(w).Encode(complaint)
}

func getUserComplaintsHandler(w http.ResponseWriter, r *http.Request) {
    secret := r.URL.Query().Get("secretCode")

    mu.Lock()
    user, ok := users[secret]
    if !ok {
        mu.Unlock()
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var result []Complaint
    for _, id := range user.Complaints {
        result = append(result, complaints[id])
    }
    mu.Unlock()

    json.NewEncoder(w).Encode(result)
}

func getAdminComplaintsHandler(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    var result []struct {
        Title string
        User  string
    }
    for _, c := range complaints {
        for _, u := range users {
            if u.ID == c.UserID {
                result = append(result, struct {
                    Title string
                    User  string
                }{Title: c.Title, User: u.Name})
                break
            }
        }
    }
    mu.Unlock()

    json.NewEncoder(w).Encode(result)
}

func viewComplaintHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    secret := r.URL.Query().Get("secretCode")

    mu.Lock()
    user, ok := users[secret]
    complaint, exists := complaints[id]
    mu.Unlock()

    if !ok || (!exists || complaint.UserID != user.ID) {
        http.Error(w, "Unauthorized or complaint not found", http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(complaint)
}

func resolveComplaintHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")

    mu.Lock()
    complaint, ok := complaints[id]
    if ok {
        complaint.Resolved = true
        complaints[id] = complaint
    }
    mu.Unlock()

    if !ok {
        http.Error(w, "Complaint not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Complaint marked resolved"})
}