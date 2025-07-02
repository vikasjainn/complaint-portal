// Common/Strvar.go
package Common

// Log Messages
const (
	LogStartingServer      = "Starting gRPC server on port "
	LogFailedToListen      = "failed to listen: %v"
	LogFailedToServe       = "failed to serve: %v"
	LogReceivedRegister    = "Received Register request for Name: %v"
	LogReceivedLogin       = "Received Login request with secret code"
	LogReceivedSubmit      = "Received SubmitComplaint request"
	LogReceivedGetUser     = "Received GetUserComplaints request"
	LogReceivedGetAdmin    = "Received GetAdminComplaints request"
	LogReceivedView        = "Received ViewComplaint request"
	LogReceivedResolve     = "Received ResolveComplaint request"
)

// Error Messages
const (
	ErrNameAndEmailRequired = "Name and Email are required"
	ErrEmailAlreadyExists   = "Email already registered"
	ErrInvalidSecretCode    = "Invalid secret code"
	ErrUnauthorized         = "Unauthorized: Invalid secret code"
	ErrComplaintNotFound    = "Complaint not found"
	ErrComplaintAccess      = "Complaint not found or you are not the owner"
)

// Success Messages
const (
	MsgComplaintResolved = "Complaint marked as resolved"
)

// Network
const (
	GRPC_Port = ":50051"
	TCP       = "tcp"
)
