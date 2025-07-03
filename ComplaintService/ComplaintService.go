// ComplaintService/ComplaintService.go
package ComplaintService

import (
	"complaint-portal/Common"
	pb "complaint-portal/Generated/ComplaintService"
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement the ComplaintServiceServer interface.
type Server struct {
	pb.UnimplementedComplaintServiceServer
}

const (
	usersCollection      = "users"
	complaintsCollection = "complaints"
)

// Register implements the Register RPC method using Firestore.
func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.User, error) {
	log.Printf(Common.LogReceivedRegister, req.GetName())

	if req.GetName() == "" || req.GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, Common.ErrNameAndEmailRequired)
	}

	// Check if email already exists by querying Firestore
	iter := Common.FirestoreClient.Collection(usersCollection).Where("Email", "==", req.GetEmail()).Limit(1).Documents(ctx)
	if _, err := iter.Next(); err != iterator.Done {
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to query database: %v", err)
		}
		return nil, status.Errorf(codes.AlreadyExists, Common.ErrEmailAlreadyExists)
	}

	user := Common.User{
		ID:         Common.GenerateID(),
		SecretCode: Common.GenerateSecretCode(),
		Name:       req.GetName(),
		Email:      req.GetEmail(),
		Complaints: []string{},
	}

	// Use the user's ID as the document ID in Firestore
	_, err := Common.FirestoreClient.Collection(usersCollection).Doc(user.ID).Set(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)
	}

	return &pb.User{
		Id:           user.ID,
		SecretCode:   user.SecretCode,
		Name:         user.Name,
		Email:        user.Email,
		ComplaintIds: user.Complaints,
	}, nil
}

// Login implements the Login RPC method using Firestore.
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.User, error) {
	log.Println(Common.LogReceivedLogin)

	iter := Common.FirestoreClient.Collection(usersCollection).Where("SecretCode", "==", req.GetSecretCode()).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, status.Errorf(codes.NotFound, Common.ErrInvalidSecretCode)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to query database: %v", err)
	}

	var user Common.User
	doc.DataTo(&user)

	return &pb.User{
		Id:           user.ID,
		SecretCode:   user.SecretCode,
		Name:         user.Name,
		Email:        user.Email,
		ComplaintIds: user.Complaints,
	}, nil
}

// SubmitComplaint implements the SubmitComplaint RPC method using Firestore.
func (s *Server) SubmitComplaint(ctx context.Context, req *pb.SubmitComplaintRequest) (*pb.Complaint, error) {
	log.Println(Common.LogReceivedSubmit)

	// Find user by secret code
	iter := Common.FirestoreClient.Collection(usersCollection).Where("SecretCode", "==", req.GetSecretCode()).Limit(1).Documents(ctx)
	userDoc, err := iter.Next()
	if err == iterator.Done {
		return nil, status.Errorf(codes.Unauthenticated, Common.ErrUnauthorized)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to query user: %v", err)
	}
	var user Common.User
	userDoc.DataTo(&user)

	// Create new complaint
	complaint := Common.Complaint{
		ID:       Common.GenerateID(),
		Title:    req.GetTitle(),
		Summary:  req.GetSummary(),
		Severity: int(req.GetSeverity()),
		UserID:   user.ID,
		Resolved: false,
	}

	// Save complaint to Firestore
	_, err = Common.FirestoreClient.Collection(complaintsCollection).Doc(complaint.ID).Set(ctx, complaint)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create complaint: %v", err)
	}

	// Update user's complaints list
	_, err = Common.FirestoreClient.Collection(usersCollection).Doc(user.ID).Update(ctx, []firestore.Update{
		{Path: "Complaints", Value: firestore.ArrayUnion(complaint.ID)},
	})
	if err != nil {
		// Attempt to roll back or log error
		return nil, status.Errorf(codes.Internal, "Failed to update user with new complaint: %v", err)
	}

	return &pb.Complaint{
		Id:       complaint.ID,
		Title:    complaint.Title,
		Summary:  complaint.Summary,
		Severity: int32(complaint.Severity),
		UserId:   complaint.UserID,
		Resolved: complaint.Resolved,
	}, nil
}

// GetUserComplaints implements the GetUserComplaints RPC method using Firestore.
func (s *Server) GetUserComplaints(ctx context.Context, req *pb.GetUserComplaintsRequest) (*pb.GetUserComplaintsResponse, error) {
	log.Println(Common.LogReceivedGetUser)

	// Find user by secret code
	iter := Common.FirestoreClient.Collection(usersCollection).Where("SecretCode", "==", req.GetSecretCode()).Limit(1).Documents(ctx)
	userDoc, err := iter.Next()
	if err == iterator.Done {
		return nil, status.Errorf(codes.Unauthenticated, Common.ErrUnauthorized)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to query user: %v", err)
	}
	var user Common.User
	userDoc.DataTo(&user)

	var result []*pb.Complaint
	// Find all complaints for that user
	complaintsIter := Common.FirestoreClient.Collection(complaintsCollection).Where("UserID", "==", user.ID).Documents(ctx)
	for {
		complaintDoc, err := complaintsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to retrieve complaints: %v", err)
		}
		var c Common.Complaint
		complaintDoc.DataTo(&c)
		result = append(result, &pb.Complaint{
			Id:       c.ID,
			Title:    c.Title,
			Summary:  c.Summary,
			Severity: int32(c.Severity),
			UserId:   c.UserID,
			Resolved: c.Resolved,
		})
	}

	return &pb.GetUserComplaintsResponse{Complaints: result}, nil
}

// GetAdminComplaints implements the GetAdminComplaints RPC method using Firestore.
func (s *Server) GetAdminComplaints(ctx context.Context, req *pb.GetAdminComplaintsRequest) (*pb.GetAdminComplaintsResponse, error) {
	log.Println(Common.LogReceivedGetAdmin)

	var result []*pb.AdminComplaintDetails
	complaintsIter := Common.FirestoreClient.Collection(complaintsCollection).Documents(ctx)
	for {
		complaintDoc, err := complaintsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to retrieve complaints: %v", err)
		}

		var c Common.Complaint
		complaintDoc.DataTo(&c)

		// Get the user for this complaint
		userDoc, err := Common.FirestoreClient.Collection(usersCollection).Doc(c.UserID).Get(ctx)
		if err != nil {
			// Log the error but continue, maybe the user was deleted
			log.Printf("Could not find user %s for complaint %s: %v", c.UserID, c.ID, err)
			continue
		}
		var u Common.User
		userDoc.DataTo(&u)

		result = append(result, &pb.AdminComplaintDetails{
			Title:    c.Title,
			UserName: u.Name,
		})
	}

	return &pb.GetAdminComplaintsResponse{Complaints: result}, nil
}

// ViewComplaint implements the ViewComplaint RPC method using Firestore with corrected logic.
func (s *Server) ViewComplaint(ctx context.Context, req *pb.ViewComplaintRequest) (*pb.Complaint, error) {
	log.Println(Common.LogReceivedView)

	// Step 1: Get the requested complaint first.
	complaintDoc, err := Common.FirestoreClient.Collection(complaintsCollection).Doc(req.GetComplaintId()).Get(ctx)
	if err != nil {
		// If the complaint doesn't exist at all, return NotFound. This is correct.
		return nil, status.Errorf(codes.NotFound, Common.ErrComplaintNotFound)
	}
	var complaint Common.Complaint
	complaintDoc.DataTo(&complaint)

	// Step 2: Now, authenticate the user making the request.
	iter := Common.FirestoreClient.Collection(usersCollection).Where("SecretCode", "==", req.GetSecretCode()).Limit(1).Documents(ctx)
	userDoc, err := iter.Next()
	if err == iterator.Done {
		// If the secret code is invalid, the user is unauthenticated.
		return nil, status.Errorf(codes.Unauthenticated, Common.ErrUnauthorized)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to query user: %v", err)
	}
	var user Common.User
	userDoc.DataTo(&user)

	// Step 3: Finally, check for ownership. This is the authorization step.
	if complaint.UserID != user.ID {
		// The user is authenticated, but not authorized to see this specific complaint.
		return nil, status.Errorf(codes.PermissionDenied, Common.ErrComplaintAccess)
	}

	// If all checks pass, return the complaint data.
	return &pb.Complaint{
		Id:       complaint.ID,
		Title:    complaint.Title,
		Summary:  complaint.Summary,
		Severity: int32(complaint.Severity),
		UserId:   complaint.UserID,
		Resolved: complaint.Resolved,
	}, nil
}


// ResolveComplaint implements the ResolveComplaint RPC method using Firestore.
func (s *Server) ResolveComplaint(ctx context.Context, req *pb.ResolveComplaintRequest) (*pb.ResolveComplaintResponse, error) {
	log.Println(Common.LogReceivedResolve)

	// Update the complaint document
	_, err := Common.FirestoreClient.Collection(complaintsCollection).Doc(req.GetComplaintId()).Update(ctx, []firestore.Update{
		{Path: "Resolved", Value: true},
	})

	if err != nil {
		// This will return a NotFound error if the document ID doesn't exist.
		return nil, status.Errorf(codes.Internal, "Failed to update complaint: %v", err)
	}

	return &pb.ResolveComplaintResponse{Message: Common.MsgComplaintResolved}, nil
}
