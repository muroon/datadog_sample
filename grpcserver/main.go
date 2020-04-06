package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/muroon/datadog_sample/config"
	"github.com/muroon/datadog_sample/grpcserver/service"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	pb "github.com/muroon/datadog_sample/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

const (
	datadogService = "grpc-server-service"
)

type server struct{}

func (s *server) GetMessages(ctx context.Context, params *pb.EmptyParams) (*pb.Messages, error) {
	ms, err := service.GetMessages(ctx)
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

func getDBTypeAndDataSource(
	conf config.Config,
) (dbType string, dbSource string, err error) {
	dbType, err = conf.GrpcDBType()
	if err != nil {
		return
	}

	dbSource, err = conf.GrpcDBDataSource()
	return
}

func main() {
	// config
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
		return
	}

	// db config
	dbType, dbSource, err := getDBTypeAndDataSource(conf)
	if err != nil {
		log.Fatalf("invalid config: %v", err)
		return
	}

	_, grpcPort, err := conf.GrpcHostAndPort()
	if err != nil {
		log.Fatalf("invalid config: %v", err)
		return
	}

	// Datadog
	tracer.Start(
		tracer.WithEnv("sample"),
	)
	defer tracer.Stop()

	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName(datadogService))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName(datadogService))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	service.InitDB(dbType)
	err = service.OpenDB(dbType, dbSource)
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
