package main

import (
    "context"
    "fmt"
    "log"

    "google.golang.org/grpc"

    "com.deer/grpclearn/pkg/gtls"
    pb "com.deer/grpclearn/proto"
)

const PORT = 9003

func main() {
    tlsClient := gtls.Client{
        ServerName: "grpcLearn",
        CertFile:   "../../conf/server/server.pem",
    }
    c, err := tlsClient.GetTLSCredentials()
    if err != nil {
        log.Fatalf("tlsClient.GetTLSCredentials err: %v", err)
    }
    // tlsClient := gtls.Client{
    //     ServerName: "grpcLearn",
    //     CertFile:   "../../conf/client/client.pem",
    //     KeyFile:    "../../conf/client/client.key",
    //     CaFile:     "../../conf/ca.pem",
    // }
    // c, err := tlsClient.GetCredentialsByCA()
    // if err != nil {
    //     log.Fatalf("tlsClient.GetCredentialsByCA err: %v", err)
    // }

    conn, err := grpc.Dial(fmt.Sprintf(":%d", PORT), grpc.WithTransportCredentials(c))
    if err != nil {
        log.Fatalf("grpc.Dial err: %v", err)
    }
    defer conn.Close()

    client := pb.NewSearchServiceClient(conn)
    resp, err := client.Search(context.Background(), &pb.SearchRequest{
        Request: "gRPC",
    })
    if err != nil {
        log.Fatalf("client.Search err: %v", err)
    }
    log.Printf("resp: %s", resp.GetResponse())
}
