package main

import (
	"grpc_datadog/server/service"
	"log"
	"net"
	"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "grpc_datadog/proto"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}


func (s *server) GetMessages(ctx context.Context, params *pb.EmptyParams) (*pb.Messages, error) {
	ms, err := service.GetMessages()
	list := make([]*pb.Message, 0, len(ms))
	for _, m := range ms {
		mes := &pb.Message{
			Id: m.Id,
			Text: m.Text,
		}
		list = append(list, mes)
	}

	return &pb.Messages{List: list}, err
}

func (s *server) PostMessage(stream pb.DataManager_PostMessageServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		m, err := service.PostMessage(in.Text)
		if err != nil {
			return err
		}

		stream.Send(&pb.Message{Id: m.Id, Text: m.Text})
	}
	return nil // RPC終了
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = service.OpenDB()
	defer service.CloseDB()
	if err != nil {
		log.Fatalf("open DB error :%v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDataManagerServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

