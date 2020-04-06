package usecases

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/grpc"

	"github.com/muroon/datadog_sample/httpserver/jsonmodel"

	pb "github.com/muroon/datadog_sample/proto"

	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

var c pb.DataManagerClient
var stream pb.DataManager_PostMessageClient
var conn *grpc.ClientConn

func openGrpc() error {
	ui := grpctrace.UnaryClientInterceptor(
		grpctrace.WithServiceName("grpc-client-service"),
	)

	// grpc connection
	var err error
	conn, err = grpc.Dial(address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(ui),
	)
	if err != nil {
		return err
	}
	c = pb.NewDataManagerClient(conn)

	// post message stream
	ctx := context.Background()
	stream, err = c.PostMessage(ctx)
	if err != nil {
		return err
	}

	go func() {
		for {
			mes, err := stream.Recv()
			if err != nil {
				log.Printf("could not greet: %v", err)
			}

			log.Printf("post message Id:%d, Text:%s", mes.Id, mes.Text)
		}
	}()
	return nil
}

func closeGrpc() error {
	err := conn.Close()
	if err != nil {
		log.Printf("close connect error: %v", err)
	}
	return err
}

func GrpcList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := c.GetMessages(ctx, new(pb.EmptyParams))
	if err != nil {
		renderErrorJSON(w, err)
		return
	}

	list := make([]*jsonmodel.Message, 0, len(res.List))
	for _, mess := range res.List {
		list = append(list, &jsonmodel.Message{ID: mess.Id, Text: mess.Text})
	}

	b, err := json.Marshal(list)
	if err != nil {
		renderErrorJSON(w, err)
		return
	}

	renderJSON(w, b)
}

func GrpcPost(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	err := stream.Send(&pb.Message{Text: text})
	if err != nil {
		renderErrorJSON(w, err)
		return
	}

	res := &jsonmodel.PostMessageResult{Status: true}
	b, err := json.Marshal(res)
	if err != nil {
		renderErrorJSON(w, err)
		return
	}

	renderJSON(w, b)
}
