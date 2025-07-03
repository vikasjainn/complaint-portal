# Complaint Portal gRPC Service

This project is a backend service for a complaint management system, built using Go, gRPC, and Google Firestore. It provides a complete set of functionalities for user registration, login, and complaint submission and management, all adhering to a strict set of coding protocols.

---

## Features

User Management: Secure user registration and login via a secret code.
Complaint Submission: Authenticated users can submit new complaints with a title, summary, and severity level.
Complaint Viewing: Users can view their own complaints, and an admin endpoint is available to view all complaints.
Complaint Resolution: An endpoint to mark complaints as resolved.
Persistent Storage: Uses Google Firestore to permanently store all user and complaint data.
Automated Testing: Includes a full suite of unit tests that run against a local Firestore emulator.

---

## Technology Stack

Language: Go
API Framework: gRPC
Data Serialization: Protocol Buffers (proto3)
Database: Google Cloud Firestore
Testing: Go's native testing package, Firestore Emulator

---

## Project Structure

The project follows a clean, protocol-driven architecture:


complaint-portal/
├── Common/                  # Shared code: models, utils, Firebase connection
├── ComplaintService/        # gRPC service implementation and test files
├── Generated/               # Auto-generated gRPC and Protobuf Go code
├── proto/                   # .proto file defining the API contract
├── test-client/             # Separate interactive CLI client
├── .gitignore               # Ensures credentials are not pushed to Git
├── go.mod                   # Go module dependencies
├── go.sum
├── main.go                  # Main application entry point
└── credentials.json         # (Ignored by Git) Firebase service account key


---

## Setup and Installation

### Prerequisites

Before you begin, ensure you have the following installed:
1.  Go: Version 1.18 or newer.
2.  Protocol Buffer Compiler (`protoc`): [Installation Guide](https://grpc.io/docs/protoc-installation/)
3.  Java JRE: Version 8 or newer (required for the Firestore emulator).
4.  Google Cloud SDK: [Installation Guide](https://cloud.google.com/sdk/docs/install)
    - After installing, add the Firestore emulator component: `gcloud components install firestore-emulator`

### Installation Steps

1.  Clone the Repository
    ```bash
    git clone <your-repository-url>
    cd complaint-portal
    ```

2.  Firebase Credentials
    - Go to your Firebase project settings, navigate to "Service accounts," and generate a new private key.
    - Rename the downloaded JSON file to `credentials.json` and place it in the root of the project directory.

3.  Install Go Plugins
    - Run these commands to install the necessary Go code generators for gRPC.
    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    ```

4.  Generate gRPC Code
    - From the project root, run the `protoc` command to generate the Go code from your `.proto` file.
    ```bash
    protoc --go_out=. --go-grpc_out=. proto/complaint.proto
    ```

5.  Install Go Dependencies
    - Tidy the modules to download all required packages.
    ```bash
    go mod tidy
    ```

---

## How to Run

### 1. Run the Server

-   From the project root (`complaint-portal/`), run the main application.
    ```bash
    go run .
    ```
-   You should see log messages indicating a successful connection to Firestore and the server starting on port `:50051`.

### 2. Run the Interactive Client

-   Open a new terminal window.
-   Navigate to the `test-client/` directory.
-   Run the client application.
    ```bash
    go run .
    ```
-   Follow the on-screen menu to register, log in, and manage complaints.

---

## How to Test

### 1. Start the Firestore Emulator

-   Open a new, dedicated terminal window.
-   Run the command to start the local Firestore emulator.
    ```bash
    gcloud emulators firestore start --host-port="localhost:8081"
    ```
-   Keep this terminal open while you run the tests.

### 2. Run the Unit Tests

-   Open another new terminal window.
-   Navigate to the project root (`complaint-portal/`).
-   Run the test suite.
    ```bash
    go test ./...
    ```
-   You should see an `ok` message for both the `Common` and `ComplaintService` packages, indicating that all tests have passed.

---

## Security

The `credentials.json` file provides full administrative access to your Firebase project and must never be committed to version control. The `.gitignore` file is configured to prevent this file from being tracked by Git.
