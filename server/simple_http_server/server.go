package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "strings"

    "google.golang.org/grpc"

    "com.deer/grpclearn/pkg/gtls"
    pb "com.deer/grpclearn/proto"
)

type SearchService struct {
    pb.UnimplementedStreamServiceServer
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
    log.Printf("Received: %v", r.GetRequest())
    return &pb.SearchResponse{Response: r.GetRequest() + " HTTP Server"}, nil
}

const PORT = 9003

func main() {
    certFile := "../../conf/server/server.pem"
    keyFile := "../../conf/server/server.key"
    tlsServer := gtls.Server{
        CertFile: certFile,
        KeyFile:  keyFile,
    }
    c, err := tlsServer.GetTLSCredentials()
    if err != nil {
        log.Fatalf("tlsServer.GetTLSCredentials err: %v", err)
    }

    // certFile := "../../conf/server/server.pem"
    // keyFile := "../../conf/server/server.key"
    // caFile := "../../conf/ca.pem"
    // tlsServer := gtls.Server{
    //     CertFile: certFile,
    //     KeyFile:  keyFile,
    //     CaFile:   caFile,
    // }
    // c, err := tlsServer.GetCredentialsByCA()
    // if err != nil {
    //     log.Fatalf("tlsServer.GetCredentialsByCA err: %v", err)
    // }

    mux := GetHTTPServeMux()
    server := grpc.NewServer(grpc.Creds(c))
    pb.RegisterSearchServiceServer(server, &SearchService{})

    _ = http.ListenAndServeTLS(fmt.Sprintf(":%d", PORT),
        certFile,
        keyFile,
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
                server.ServeHTTP(w, r)
            } else {
                mux.ServeHTTP(w, r)
            }
            return
        }),
    )
}

func GetHTTPServeMux() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Printf(r.RemoteAddr)
        _, _ = w.Write([]byte("zhangjinfu: go-grpc-example"))
    })
    return mux
}
