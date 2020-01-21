package main

import (
    "fmt"
    "io"
    "log"
    "net"

    "google.golang.org/grpc"

    pb "com.deer/grpclearn/proto"
)

const (
    PORT = 9002
)

func main() {
    server := grpc.NewServer()
    pb.RegisterStreamServiceServer(server, &StreamService{})

    listen, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    if err = server.Serve(listen); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

type StreamService struct {
    pb.UnimplementedStreamServiceServer
}

func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
    for n := 0; n < 6; n++ {
        err := stream.Send(&pb.StreamResponse{
            Pt: &pb.StreamPoint{
                Name:  r.Pt.GetName(),
                Value: r.Pt.GetValue() + int32(n),
            },
        })
        if err != nil {
            return err
        }
    }
    return nil
}

func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
    for {
        r, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.StreamResponse{Pt: &pb.StreamPoint{Name: "gPRC Stream Server: Record", Value: 1}})
        }
        if err != nil {
            return err
        }

        log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.GetPt().GetName(), r.GetPt().GetValue())
    }
}

func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
    n := 0
    for {
        err := stream.Send(&pb.StreamResponse{
            Pt: &pb.StreamPoint{
                Name:  "gPRC Stream Server: Route",
                Value: int32(n),
            },
        })
        if err != nil {
            return err
        }
        r, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        n++
        log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.GetPt().GetName(), r.GetPt().GetValue())
    }
}
