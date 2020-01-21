package main

import (
    "context"
    "fmt"
    "io"
    "log"

    "google.golang.org/grpc"

    pb "com.deer/grpclearn/proto"
)

const (
    PORT = 9002
)

func main() {
    conn, err := grpc.Dial(fmt.Sprintf(":%d", PORT), grpc.WithInsecure())
    if err != nil {
        log.Fatalf("failed to dial: %v", err)
    }
    defer conn.Close()

    client := pb.NewStreamServiceClient(conn)
    err = printLists(client, &pb.StreamRequest{
        Pt: &pb.StreamPoint{
            Name:  "gRPC Stream Client: List",
            Value: 2020,
        },
    })
    if err != nil {
        log.Fatalf("failed to printLists: %v", err)
    }
    err = printRecord(client, &pb.StreamRequest{
        Pt: &pb.StreamPoint{
            Name:  "gRPC Stream Client: Record",
            Value: 2020,
        },
    })
    if err != nil {
        log.Fatalf("failed to printRecord: %v", err)
    }
    err = printRoute(client, &pb.StreamRequest{
        Pt: &pb.StreamPoint{
            Name:  "gRPC Stream Client: Route",
            Value: 2020,
        },
    })
    if err != nil {
        log.Fatalf("failed to printRoute: %v", err)
    }
}

func printLists(client pb.StreamServiceClient, r *pb.StreamRequest) error {
    stream, err := client.List(context.Background(), r)
    if err != nil {
        return err
    }
    for {
        resp, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }

        log.Printf("resp: pj.name: %s, pj.value: %d", resp.GetPt().GetName(), resp.GetPt().GetValue())
    }
    return nil
}

func printRecord(client pb.StreamServiceClient, r *pb.StreamRequest) error {
    stream, err := client.Record(context.Background())
    if err != nil {
        return err
    }
    for n := 0; n < 6; n++ {
        err := stream.Send(r)
        if err != nil {
            return err
        }
    }

    resp, err := stream.CloseAndRecv()
    if err != nil {
        return err
    }
    log.Printf("resp: pj.name: %s, pj.value: %d", resp.GetPt().GetName(), resp.GetPt().GetValue())

    return nil
}

func printRoute(client pb.StreamServiceClient, r *pb.StreamRequest) error {
    stream, err := client.Route(context.Background())
    if err != nil {
        return err
    }
    for n := 0; n <= 6; n++ {
        err = stream.Send(r)
        if err != nil {
            return err
        }
        resp, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        log.Printf("resp: pj.name: %s, pt.value: %d", resp.GetPt().GetName(), resp.GetPt().GetValue())
    }
    _ = stream.CloseSend()
    return nil
}
