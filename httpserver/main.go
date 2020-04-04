package main

import (
	"datadog_sample/httpserver/usecases"
	"log"
	"net/http"
)


func main() {
	err := usecases.Init()
	defer usecases.End()
	if err != nil {
		log.Fatalf("init error: %v", err)
		return
	}

	// http handle
	http.HandleFunc("/grpc/", usecases.GrpcList)
	http.HandleFunc("/grpc/post", usecases.GrpcPost)
	http.HandleFunc("/custom/", usecases.Custom)
	http.HandleFunc("/db/", usecases.DBList)
	http.HandleFunc("/db/post", usecases.DBPost)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("http listenServe error: %v", err)
	}
}

