package main

import (
	"context"
	"encoding/json"
	"grpc_datadog/httpserver/jsonmodel"
	"log"
	"net/http"

	"google.golang.org/grpc"
	pb "grpc_datadog/proto"
)

const(
	address     = "localhost:50051"
)

var c pb.DataManagerClient
var stream pb.DataManager_PostMessageClient

func main() {
	// grpc connection
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c = pb.NewDataManagerClient(conn)

	// post message stream
	ctx := context.Background()
	stream, err = c.PostMessage(ctx)
	if err != nil {
		log.Fatalf("could not stream: %v", err)
		return
	}

	go func() {
		for {
			mes, err := stream.Recv()
			if err != nil {
				log.Fatalf("could not greet: %v", err)
				break
			}

			log.Printf("post message Id:%d, Text:%s", mes.Id, mes.Text)
		}
	}()

	// http handle
	http.HandleFunc("/", list)
	http.HandleFunc("/post", post)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("http listenServe error: %v", err)
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := c.GetMessages(ctx, new(pb.EmptyParams))
	if err != nil {
		log.Fatalf("could not unary message: %v", err)
		return
	}

	list := make([]*jsonmodel.Message, 0, len(res.List))
	for _, mess := range res.List {
		list = append(list, &jsonmodel.Message{ID: mess.Id, Text: mess.Text})
	}

	b, err := json.Marshal(list)
	if err != nil {
		log.Fatalf("json marshal error: %v", err)
		return
	}

	renderJSON(w, b)
}

func post(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	err := stream.Send(&pb.Message{Text: text})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}

	res := &jsonmodel.PostMessageResult{Status:true}
	b, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("json marshal error: %v", err)
		return
	}

	renderJSON(w, b)
}

func renderJSON(w http.ResponseWriter, b []byte) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}
