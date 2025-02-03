package gapi

import "f/pb"

// Server serves gRPC requests.
type Server struct {
	pb.UnimplementedHelloServiceServer
}

func NewServer() *Server {
	return &Server{}
}


