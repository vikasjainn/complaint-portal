// ComplaintService/ComplaintService.go
// This file contains the API logic for the ComplaintService.
package ComplaintService

import (
	"complaint-portal/Common"
	pb "complaint-portal/Generated/ComplaintService" // Import generated protobuf package
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement the ComplaintServiceServer interface.
type Server struct {
	pb.UnimplementedComplaintServiceServer
}

// Register implements the Register RPC method.
func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.User, error) {
	log.Printf(Common.LogReceivedRegister, req.GetName())

	if req.GetName() == "" || req.GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, Common.ErrNameAndEmailRequired)
	}

	Common.Mu.Lock()
	defer Common.Mu.Unlock()

	for _, u := range Common.Users {
		if u.Email == req.GetEmail() {
			return nil, status.Errorf(codes.AlreadyExists, Common.ErrEmailAlreadyExists)
		}
	}

	user := Common.User{
		ID:         Common.GenerateID(),
		SecretCode: Common.GenerateSecretCode(),
		Name:       req.GetName(),
		Email:      req.GetEmail(),
		Complaints: []string{},
	}
	Common.Users[user.SecretCode] = user

	return &pb.User{
		Id:           user.ID,
		SecretCode:   user.SecretCode,
		Name:         user.Name,
		Email:        user.Email,
		ComplaintIds: user.Complaints,
	}, nil
}

// Login implements the Login RPC method.
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.User, error) {
	log.Println(Common.LogReceivedLogin)

	Common.Mu.Lock()
	user, ok := Common.Users[req.GetSecretCode()]
	Common.Mu.Unlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, Common.ErrInvalidSecretCode)
	}

	return &pb.User{
		Id:           user.ID,
		SecretCode:   user.SecretCode,
		Name:         user.Name,
		Email:        user.Email,
		ComplaintIds: user.Complaints,
	}, nil
}

// SubmitComplaint implements the SubmitComplaint RPC method.
func (s *Server) SubmitComplaint(ctx context.Context, req *pb.SubmitComplaintRequest) (*pb.Complaint, error) {
	log.Println(Common.LogReceivedSubmit)

	Common.Mu.Lock()
	defer Common.Mu.Unlock()

	user, ok := Common.Users[req.GetSecretCode()]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, Common.ErrUnauthorized)
	}

	complaint := Common.Complaint{
		ID:       Common.GenerateID(),
		Title:    req.GetTitle(),
		Summary:  req.GetSummary(),
		Severity: int(req.GetSeverity()),
		UserID:   user.ID,
		Resolved: false,
	}

	Common.Complaints[complaint.ID] = complaint
	user.Complaints = append(user.Complaints, complaint.ID)
	Common.Users[req.GetSecretCode()] = user

	return &pb.Complaint{
		Id:       complaint.ID,
		Title:    complaint.Title,
		Summary:  complaint.Summary,
		Severity: int32(complaint.Severity),
		UserId:   complaint.UserID,
		Resolved: complaint.Resolved,
	}, nil
}

// GetUserComplaints implements the GetUserComplaints RPC method.
func (s *Server) GetUserComplaints(ctx context.Context, req *pb.GetUserComplaintsRequest) (*pb.GetUserComplaintsResponse, error) {
	log.Println(Common.LogReceivedGetUser)

	Common.Mu.Lock()
	defer Common.Mu.Unlock()

	user, ok := Common.Users[req.GetSecretCode()]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, Common.ErrUnauthorized)
	}

	var result []*pb.Complaint
	for _, id := range user.Complaints {
		if c, exists := Common.Complaints[id]; exists {
			result = append(result, &pb.Complaint{
				Id:       c.ID,
				Title:    c.Title,
				Summary:  c.Summary,
				Severity: int32(c.Severity),
				UserId:   c.UserID,
				Resolved: c.Resolved,
			})
		}
	}

	return &pb.GetUserComplaintsResponse{Complaints: result}, nil
}

// GetAdminComplaints implements the GetAdminComplaints RPC method.
func (s *Server) GetAdminComplaints(ctx context.Context, req *pb.GetAdminComplaintsRequest) (*pb.GetAdminComplaintsResponse, error) {
	log.Println(Common.LogReceivedGetAdmin)

	Common.Mu.Lock()
	defer Common.Mu.Unlock()

	var result []*pb.AdminComplaintDetails
	// This nested loop is inefficient but matches the original project's logic.
	for _, c := range Common.Complaints {
		for _, u := range Common.Users {
			if u.ID == c.UserID {
				result = append(result, &pb.AdminComplaintDetails{
					Title:    c.Title,
					UserName: u.Name,
				})
				break
			}
		}
	}

	return &pb.GetAdminComplaintsResponse{Complaints: result}, nil
}

// ViewComplaint implements the ViewComplaint RPC method.
func (s *Server) ViewComplaint(ctx context.Context, req *pb.ViewComplaintRequest) (*pb.Complaint, error) {
	log.Println(Common.LogReceivedView)

	Common.Mu.Lock()
	defer Common.Mu.Unlock()

	user, ok := Common.Users[req.GetSecretCode()]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, Common.ErrUnauthorized)
	}

	complaint, exists := Common.Complaints[req.GetComplaintId()]
	if !exists || complaint.UserID != user.ID {
		return nil, status.Errorf(codes.NotFound, Common.ErrComplaintAccess)
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

// ResolveComplaint implements the ResolveComplaint RPC method.
func (s *Server) ResolveComplaint(ctx context.Context, req *pb.ResolveComplaintRequest) (*pb.ResolveComplaintResponse, error) {
	log.Println(Common.LogReceivedResolve)

	Common.Mu.Lock()
	defer Common.Mu.Unlock()

	complaint, ok := Common.Complaints[req.GetComplaintId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, Common.ErrComplaintNotFound)
	}

	complaint.Resolved = true
	Common.Complaints[req.GetComplaintId()] = complaint

	return &pb.ResolveComplaintResponse{Message: Common.MsgComplaintResolved}, nil
}
