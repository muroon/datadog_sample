package main

import (
	"github.com/muroon/datadog_sample/grpcserver/service"
	"io"
	"log"
	"net"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	pb "github.com/muroon/datadog_sample/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

const (
	port           = ":50051"
	datadogService = "grpc-server-service"
)

type server struct{}

func (s *server) GetMessages(ctx context.Context, params *pb.EmptyParams) (*pb.Messages, error) {
	ms, err := service.GetMessages()
	list := make([]*pb.Message, 0, len(ms))
	for _, m := range ms {
		mes := &pb.Message{
			Id:   m.Id,
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
	// Datadog
	tracer.Start(
		tracer.WithEnv("sample"),
	)
	defer tracer.Stop()

	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName(datadogService))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName(datadogService))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	err = service.OpenDB()
	defer service.CloseDB()
	if err != nil {
		log.Fatalf("open DB error :%v", err)
		return
	}

	s := grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))
	pb.RegisterDataManagerServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
}
