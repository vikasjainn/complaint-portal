## Complaint Portal Backend API

This is a small backend API built in Golang. It allows users to register, log in using a secret code, and submit complaints. There's also an option for admins to see and resolve all complaints.

## Endpoints

- `POST /register` – register a user with name and email
- `GET /login?secretCode=` – login using secretCode
- `POST /submitComplaint?secretCode=` – submit a complaint
- `GET /getAllComplaintsForUser?secretCode=` – see your own complaints
- `GET /getAllComplaintsForAdmin` – admin can view all complaints
- `GET /viewComplaint?id=&secretCode=` – view full complaint detail
- `GET /resolveComplaint?id=` – mark a complaint as resolved

## How to Run

- Install Go
- In terminal: `go run main.go handlers.go models.go store.go utils.go`
- Use Postman to test the routes

## Notes

- No database, everything is stored in memory
- Duplicate emails are blocked
- No third-party libraries used
