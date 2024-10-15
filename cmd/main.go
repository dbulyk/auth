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
	address = "localhost:8001"
	authID  = 12
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	r, err := c.Get(ctx, &desc.GetRequest{Id: authID})
	if err != nil {
		log.Panicf("failed to get note by id: %v", err)
	}

	log.Printf(color.RedString("Note info:\n"), color.GreenString("%+v", r.GetName()))
}
