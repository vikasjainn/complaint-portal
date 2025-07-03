// ComplaintService/ComplaintService_test.go
package ComplaintService

import (
	"complaint-portal/Common"
	pb "complaint-portal/Generated/ComplaintService"
	"context"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// NOTE: The collection name constants have been removed from this file
// as they are already defined in ComplaintService.go and are accessible
// within the same package.

var firestoreClient *firestore.Client

// TestMain sets up the connection to the Firestore emulator before running tests
// and closes the connection after all tests have run.
func TestMain(m *testing.M) {
	// Set the FIRESTORE_EMULATOR_HOST environment variable to point to the emulator.
	os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8081")

	ctx := context.Background()
	// The project ID for the emulator can be any string.
	client, err := firestore.NewClient(ctx, "test-project", option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	if err != nil {
		log.Fatalf("Failed to create Firestore client for emulator: %v", err)
	}
	firestoreClient = client
	Common.FirestoreClient = client // Set the global client for the service to use

	// Run the tests
	exitCode := m.Run()

	// Clean up and exit
	firestoreClient.Close()
	os.Exit(exitCode)
}

// clearFirestore deletes all documents from the specified collections.
func clearFirestore(ctx context.Context, t *testing.T) {
	collections := []string{usersCollection, complaintsCollection}
	for _, coll := range collections {
		docs, err := firestoreClient.Collection(coll).Documents(ctx).GetAll()
		if err != nil {
			t.Fatalf("Failed to get documents from %s: %v", coll, err)
		}
		for _, doc := range docs {
			_, err := doc.Ref.Delete(ctx)
			if err != nil {
				t.Fatalf("Failed to delete document %s from %s: %v", doc.Ref.ID, coll, err)
			}
		}
	}
}

// TestRegister tests the Register RPC method.
func TestRegister(t *testing.T) {
	ctx := context.Background()
	s := Server{}

	// Clean the database before each test run
	clearFirestore(ctx, t)

	// Test case 1: Successful registration
	req1 := &pb.RegisterRequest{Name: "Test User", Email: "test@example.com"}
	res1, err := s.Register(ctx, req1)
	if err != nil {
		t.Fatalf("Expected no error for successful registration, but got: %v", err)
	}
	if res1.GetName() != req1.Name || res1.GetEmail() != req1.Email {
		t.Errorf("Expected user name and email to be '%s' and '%s', got '%s' and '%s'", req1.Name, req1.Email, res1.Name, res1.Email)
	}
	if res1.GetId() == "" || res1.GetSecretCode() == "" {
		t.Error("Expected user ID and secret code to be generated")
	}

	// Test case 2: Attempt to register with the same email
	_, err = s.Register(ctx, req1)
	if err == nil {
		t.Fatal("Expected an error for duplicate email registration, but got none")
	}
	if status.Code(err) != codes.AlreadyExists {
		t.Errorf("Expected error code %v, but got %v", codes.AlreadyExists, status.Code(err))
	}

	// Test case 3: Registration with missing name
	req3 := &pb.RegisterRequest{Name: "", Email: "test2@example.com"}
	_, err = s.Register(ctx, req3)
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument error for missing name, but got %v", status.Code(err))
	}
}

// TestLogin tests the Login RPC method.
func TestLogin(t *testing.T) {
	ctx := context.Background()
	s := Server{}
	clearFirestore(ctx, t)

	// First, register a user to test login
	regReq := &pb.RegisterRequest{Name: "Login User", Email: "login@example.com"}
	regRes, _ := s.Register(ctx, regReq)

	// Test case 1: Successful login with correct secret code
	loginReq1 := &pb.LoginRequest{SecretCode: regRes.GetSecretCode()}
	loginRes1, err := s.Login(ctx, loginReq1)
	if err != nil {
		t.Fatalf("Expected no error for successful login, but got: %v", err)
	}
	if loginRes1.GetName() != regReq.Name {
		t.Errorf("Expected user name to be '%s', got '%s'", regReq.Name, loginRes1.Name)
	}

	// Test case 2: Failed login with incorrect secret code
	loginReq2 := &pb.LoginRequest{SecretCode: "invalid-secret-code"}
	_, err = s.Login(ctx, loginReq2)
	if status.Code(err) != codes.NotFound {
		t.Errorf("Expected NotFound error for invalid secret, but got %v", status.Code(err))
	}
}

// TestSubmitComplaint tests the SubmitComplaint RPC method.
func TestSubmitComplaint(t *testing.T) {
	ctx := context.Background()
	s := Server{}
	clearFirestore(ctx, t)

	// Register a user first
	regRes, _ := s.Register(ctx, &pb.RegisterRequest{Name: "Complaint Filer", Email: "filer@example.com"})

	// Test case 1: Successful complaint submission
	submitReq1 := &pb.SubmitComplaintRequest{
		SecretCode: regRes.GetSecretCode(),
		Title:      "Test Complaint",
		Summary:    "This is a test summary.",
		Severity:   3,
	}
	submitRes1, err := s.SubmitComplaint(ctx, submitReq1)
	if err != nil {
		t.Fatalf("Expected no error for successful complaint submission, but got: %v", err)
	}
	if submitRes1.GetTitle() != submitReq1.Title {
		t.Errorf("Expected complaint title to be '%s', got '%s'", submitReq1.Title, submitRes1.Title)
	}

	// Test case 2: Submission with invalid secret code
	submitReq2 := &pb.SubmitComplaintRequest{
		SecretCode: "invalid-secret",
		Title:      "Another Complaint",
		Summary:    "Summary here.",
		Severity:   1,
	}
	_, err = s.SubmitComplaint(ctx, submitReq2)
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("Expected Unauthenticated error for invalid secret, but got %v", status.Code(err))
	}
}

// TestGetUserComplaints tests the GetUserComplaints RPC method.
func TestGetUserComplaints(t *testing.T) {
	ctx := context.Background()
	s := Server{}
	clearFirestore(ctx, t)

	// Setup: Register a user and submit two complaints
	regRes, _ := s.Register(ctx, &pb.RegisterRequest{Name: "Multi Complaint User", Email: "multi@example.com"})
	_, _ = s.SubmitComplaint(ctx, &pb.SubmitComplaintRequest{SecretCode: regRes.GetSecretCode(), Title: "Complaint A"})
	_, _ = s.SubmitComplaint(ctx, &pb.SubmitComplaintRequest{SecretCode: regRes.GetSecretCode(), Title: "Complaint B"})

	// Test case 1: Get complaints for the user
	getReq := &pb.GetUserComplaintsRequest{SecretCode: regRes.GetSecretCode()}
	getRes, err := s.GetUserComplaints(ctx, getReq)
	if err != nil {
		t.Fatalf("Expected no error when getting user complaints, but got: %v", err)
	}
	if len(getRes.GetComplaints()) != 2 {
		t.Errorf("Expected to get 2 complaints, but got %d", len(getRes.GetComplaints()))
	}
}

// TestViewComplaint tests the ViewComplaint RPC method.
func TestViewComplaint(t *testing.T) {
	ctx := context.Background()
	s := Server{}
	clearFirestore(ctx, t)

	// Setup: Register two users, one submits a complaint
	user1Res, _ := s.Register(ctx, &pb.RegisterRequest{Name: "User One", Email: "one@example.com"})
	user2Res, _ := s.Register(ctx, &pb.RegisterRequest{Name: "User Two", Email: "two@example.com"})
	complaintRes, _ := s.SubmitComplaint(ctx, &pb.SubmitComplaintRequest{SecretCode: user1Res.GetSecretCode(), Title: "User One's Complaint"})

	// Test case 1: Owner tries to view their own complaint
	viewReq1 := &pb.ViewComplaintRequest{SecretCode: user1Res.GetSecretCode(), ComplaintId: complaintRes.GetId()}
	_, err := s.ViewComplaint(ctx, viewReq1)
	if err != nil {
		t.Fatalf("Expected no error for owner viewing complaint, but got: %v", err)
	}

	// Test case 2: Another user tries to view the complaint
	viewReq2 := &pb.ViewComplaintRequest{SecretCode: user2Res.GetSecretCode(), ComplaintId: complaintRes.GetId()}
	_, err = s.ViewComplaint(ctx, viewReq2)
	if status.Code(err) != codes.PermissionDenied {
		t.Errorf("Expected PermissionDenied error for non-owner viewing, but got %v", status.Code(err))
	}
}

// TestResolveComplaint tests the ResolveComplaint RPC method.
func TestResolveComplaint(t *testing.T) {
	ctx := context.Background()
	s := Server{}
	clearFirestore(ctx, t)

	// Setup: Register a user and submit a complaint
	regRes, _ := s.Register(ctx, &pb.RegisterRequest{Name: "Resolver User", Email: "resolver@example.com"})
	complaintRes, _ := s.SubmitComplaint(ctx, &pb.SubmitComplaintRequest{SecretCode: regRes.GetSecretCode(), Title: "To Be Resolved"})

	// Test case 1: Resolve the complaint
	resolveReq := &pb.ResolveComplaintRequest{ComplaintId: complaintRes.GetId()}
	_, err := s.ResolveComplaint(ctx, resolveReq)
	if err != nil {
		t.Fatalf("Expected no error when resolving complaint, but got: %v", err)
	}

	// Verify the complaint is marked as resolved
	viewReq := &pb.ViewComplaintRequest{SecretCode: regRes.GetSecretCode(), ComplaintId: complaintRes.GetId()}
	viewRes, _ := s.ViewComplaint(ctx, viewReq)
	if !viewRes.GetResolved() {
		t.Error("Expected complaint to be marked as resolved, but it was not")
	}
}

// TestGetAdminComplaints tests the GetAdminComplaints RPC method.
func TestGetAdminComplaints(t *testing.T) {
	ctx := context.Background()
	s := Server{}
	clearFirestore(ctx, t)

	// Setup: Register two users and have them each submit a complaint.
	user1Res, _ := s.Register(ctx, &pb.RegisterRequest{Name: "Admin Test User 1", Email: "admin1@example.com"})
	user2Res, _ := s.Register(ctx, &pb.RegisterRequest{Name: "Admin Test User 2", Email: "admin2@example.com"})
	_, _ = s.SubmitComplaint(ctx, &pb.SubmitComplaintRequest{SecretCode: user1Res.GetSecretCode(), Title: "Admin Complaint 1"})
	_, _ = s.SubmitComplaint(ctx, &pb.SubmitComplaintRequest{SecretCode: user2Res.GetSecretCode(), Title: "Admin Complaint 2"})

	// Test case 1: Call the admin endpoint
	adminReq := &pb.GetAdminComplaintsRequest{}
	adminRes, err := s.GetAdminComplaints(ctx, adminReq)
	if err != nil {
		t.Fatalf("Expected no error for GetAdminComplaints, but got: %v", err)
	}

	// Verify that both complaints are returned
	if len(adminRes.GetComplaints()) != 2 {
		t.Errorf("Expected admin to see 2 complaints, but got %d", len(adminRes.GetComplaints()))
	}
}
