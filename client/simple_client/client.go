package main

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"

    pb "com.deer/grpclearn/proto"
)

const PORT = 9001

func main() {
    // c, err := credentials.NewClientTLSFromFile("../../conf/server.pem", "grpcLearn")
    // if err != nil {
    //     log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
    // }
    cert, err := tls.LoadX509KeyPair("../../conf/client/client.pem", "../../conf/client/client.key")
    if err != nil {
        log.Fatalf("tls.LoadX509KeyPair err: %v", err)
    }
    certPool := x509.NewCertPool()
    ca, err := ioutil.ReadFile("../../conf/ca.pem")
    if err != nil {
        log.Fatalf("ioutil.ReadFile err: %v", err)
    }
    if ok := certPool.AppendCertsFromPEM(ca); !ok {
        log.Fatalf("certPool.AppendCertsFromPEM err")
    }
    c := credentials.NewTLS(&tls.Config{
        Certificates: []tls.Certificate{cert},
        ServerName:   "grpcLearn",
        RootCAs:      certPool,
    })

    conn, err := grpc.Dial(fmt.Sprintf(":%d", PORT), grpc.WithTransportCredentials(c))
    if err != nil {
        log.Fatalf("failed to dial: %v", err)
    }
    defer conn.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client := pb.NewSearchServiceClient(conn)
    resp, err := client.Search(ctx, &pb.SearchRequest{Request: "gRPC"})
    if err != nil {
        log.Fatalf("failed to call Search: %v", err)
    }
    log.Printf("resp: %s", resp.GetResponse())
}
