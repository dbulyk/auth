package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	desc "auth/pkg/auth_v1"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			slog.Error("conn isn't closed", "error", err.Error())
		}
	}(conn)

	c := desc.NewAuthV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: 12})
	if err != nil {
		log.Panicf("failed to get auth info by id: %v", err)
	}

	log.Printf(color.RedString("Auth info:\n"), color.GreenString("%+v", r.GetName()))
}
