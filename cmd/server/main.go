package main

import (
	desc "auth/pkg/auth_v1"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	desc.UnimplementedAuthV1Server
}

func (s *server) Create(_ context.Context, in *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Get Create req, %v", in)
	return &desc.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) Get(_ context.Context, in *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Get Get req, %v", in)
	return &desc.GetResponse{
		Id:        gofakeit.Int64(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      1,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}, nil
}

func (s *server) Update(_ context.Context, in *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("Get Update req, %v", in)
	return nil, nil
}

func (s *server) Delete(_ context.Context, in *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Get Delete req, %v", in)
	return nil, nil
}

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
