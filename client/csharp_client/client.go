package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "com.deer/grpclearn/proto"
)

const PORT = 50051

func main() {

	conn, err := grpc.Dial(fmt.Sprintf(":%d", PORT), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := pb.NewGreeterClient(conn)
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "gRPC"})
	if err != nil {
		log.Fatalf("failed to call SayHello: %v", err)
	}
	log.Printf("resp: %s", resp.GetMessage())
}
