package main

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "runtime/debug"

    "github.com/grpc-ecosystem/go-grpc-middleware"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/status"

    pb "com.deer/grpclearn/proto"
)

type SearchService struct {
    pb.UnimplementedSearchServiceServer
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
    log.Printf("Received: %v", r.GetRequest())
    return &pb.SearchResponse{
        Response: r.GetRequest() + " Server",
    }, nil
}

const PORT = 9001

func main() {
    // c, err := credentials.NewServerTLSFromFile("../../conf/server.pem", "../../conf/server.key")
    // if err != nil {
    //     log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
    // }
    cert, err := tls.LoadX509KeyPair("../../conf/server/server.pem", "../../conf/server/server.key")
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
        ClientAuth:   tls.RequireAndVerifyClientCert,
        ClientCAs:    certPool,
    })

    server := grpc.NewServer(grpc.Creds(c), grpc_middleware.WithUnaryServerChain(
        RecoveryInterceptor,
        LoggingInterceptor,
    ))
    pb.RegisterSearchServiceServer(server, &SearchService{})

    listen, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    if err = server.Serve(listen); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    log.Printf("gRPC method: %s, %v", info.FullMethod, req)
    resp, err := handler(ctx, req)
    log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
    return resp, err
}

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
    defer func() {
        if e := recover(); e != nil {
            debug.PrintStack()
            err = status.Errorf(codes.Internal, "Panic err: %v", e)
        }
    }()
    return handler(ctx, req)
}
