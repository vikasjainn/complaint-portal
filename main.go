// main.go
package main

import (
	"complaint-portal/Common"
	"complaint-portal/ComplaintService"
	pb "complaint-portal/Generated/ComplaintService"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	log.Printf(Common.LogStartingServer + Common.GRPC_Port)

	lis, err := net.Listen(Common.TCP, Common.GRPC_Port)
	if err != nil {
		log.Fatalf(Common.LogFailedToListen, err)
	}

	s := grpc.NewServer()

	// Register our server implementation
	pb.RegisterComplaintServiceServer(s, &ComplaintService.Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf(Common.LogFailedToServe, err)
	}
}
