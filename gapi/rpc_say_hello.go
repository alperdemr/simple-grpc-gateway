package gapi

import (
	"context"
	"f/pb"
)

func (server *Server) SayHello( ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello, " + req.GetName()}, nil
	
}