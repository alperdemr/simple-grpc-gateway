package main

import (
	"context"
	"f/api"
	"f/gapi"
	"f/pb"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main () {
	go runGrpcServer()
	runGatewayServer()
}

func runFiberServer() {
	server := api.NewServer()
	server.Start(":8080")
}

func runGrpcServer() {
	grpcServer := grpc.NewServer()
	server := gapi.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener,err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal("cannot create listener")
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}

	
}

func runGatewayServer() {
	server := gapi.NewServer()
	grpcMux := runtime.NewServeMux()
	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()
	err := pb.RegisterHelloServiceHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listener,err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("cannot create listener", err)
	}

	log.Printf("start gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start gateway server",err)
	}




	
}