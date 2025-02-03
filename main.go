package main

import (
	"context"
	"f/api"
	"f/gapi"
	"f/pb"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main () {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	go runGrpcServer()
	runGatewayServer()
}

func runFiberServer() {
	server := api.NewServer()
	server.Start(":8080")
}

func runGrpcServer() {
	
	server := gapi.NewServer()

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterHelloServiceServer(grpcServer, server)
	reflection.Register(grpcServer)



	listener,err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}

	
}

func runGatewayServer() {
	server := gapi.NewServer()
	grpcMux := runtime.NewServeMux()
	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()
	err := pb.RegisterHelloServiceHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msg("cannot register gateway server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listener,err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	log.Info().Msgf("start gateway server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg("cannot start gateway server")

	}
}