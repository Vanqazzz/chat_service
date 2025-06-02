package main

import (
	"context"
	"fmt"
	"log"
	"time"

	protos "github.com/Vanqazzz/protos/gen/go/chat_service/chat"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:44044", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could to connect:", err)
	}
	defer conn.Close()

	client := protos.NewChatServiceClient(conn)

	stream, err := client.ReceiveMessages(context.Background(), &protos.UserStreamRequest{
		UserId: "user1",
	})
	if err != nil {
		log.Fatalf("errror on ReceiveMessage: %v", err)
	}

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Fatalf("stream error: %v", err)
			}
			fmt.Printf("Received message: %+v\n", msg)
		}
	}()

	time.Sleep(5 * time.Minute)
}
