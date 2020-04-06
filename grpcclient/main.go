package main

import (
	"context"
	"log"
	"os"
	"time"
	"fmt"

	"google.golang.org/grpc"
	pb "github.com/muroon/datadog_sample/proto"
)

const (
	address     = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewDataManagerClient(conn)


	content := "test"
	if len(os.Args) > 1 {
		content = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// stream
	stream, err := c.PostMessage(ctx)
	if err != nil {
		log.Fatalf("could not stream: %v", err)
		return
	}

	var end bool

	go func() {
		var cnt int
		for {
			response, err := stream.Recv()
			if err != nil {
				log.Fatalf("could not greet: %v", err)
				end = true
				break
			}

			cnt++

			log.Printf("Id:%d, Text:%s", response.Id, response.Text)

			if cnt == 5 {
				end = true
				break
			}
		}
	}()

	for i := 0; i < 5; i++ {
		err = stream.Send(&pb.Message{Text: fmt.Sprintf("stream message %s %d", content, i)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
			break
		}
	}

	for {
		if end {
			break
		}
	}

	r, err := c.GetMessages(ctx, new(pb.EmptyParams))
	if err != nil {
		log.Fatalf("could not unary message: %v", err)
		return
	}

	for _, mess := range r.List {
		log.Printf("GetMessages. Id:%d, Text:%s", mess.Id, mess.Text)
	}
}



